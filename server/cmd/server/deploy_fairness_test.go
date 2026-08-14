package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"server/internal/room"
	"server/internal/ws"
)

type deployQuestionProvider struct {
	question room.Question
}

func (p deployQuestionProvider) AssignQuestion(string) (room.Question, error) {
	return p.question, nil
}

func setupDeployServer(t *testing.T, provider room.QuestionProvider) (*ws.Hub, *httptest.Server) {
	t.Helper()
	hub := ws.NewHub(provider)
	go hub.Run()
	server := httptest.NewServer(buildDeployMux(hub, nil, nil))
	t.Cleanup(func() {
		server.Close()
		hub.Stop()
	})
	return hub, server
}

func dialDeployWS(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial /ws via deploy mux: %v", err)
	}
	return conn
}

func deploySendJoin(t *testing.T, conn *websocket.Conn, name, roomID string) {
	t.Helper()
	payload, _ := json.Marshal(ws.JoinPayload{Name: name, RoomID: roomID})
	if err := conn.WriteJSON(ws.Message{Type: ws.MessageTypeJoin, RoomID: roomID, Payload: payload}); err != nil {
		t.Fatalf("join: %v", err)
	}
}

func deployReadUntil(t *testing.T, conn *websocket.Conn, want string) ([]byte, ws.Message) {
	t.Helper()
	for i := 0; i < 20; i++ {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read waiting for %q: %v", want, err)
		}
		var msg ws.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if msg.Type == want {
			return raw, msg
		}
	}
	t.Fatalf("did not receive %q", want)
	return nil, ws.Message{}
}

func TestDeployMux_HealthAndEmbeddedAssets(t *testing.T) {
	hub := ws.NewHub(nil)
	go hub.Run()
	defer hub.Stop()

	server := httptest.NewServer(buildDeployMux(hub, nil, nil))
	defer server.Close()

	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || string(body) != "ok" {
		t.Fatalf("health status=%d body=%q", res.StatusCode, body)
	}

	indexRes, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	defer indexRes.Body.Close()
	indexBody, _ := io.ReadAll(indexRes.Body)
	if indexRes.StatusCode != http.StatusOK || !bytes.Contains(indexBody, []byte("<html")) {
		t.Fatalf("expected embedded index.html, status=%d", indexRes.StatusCode)
	}

	sub, err := fs.Sub(embeddedDist, "dist")
	if err != nil {
		t.Fatalf("embed subtree: %v", err)
	}
	html := string(indexBody)
	for _, prefix := range []string{`src="/assets/`, `href="/assets/`} {
		rest := html
		for {
			i := strings.Index(rest, prefix)
			if i < 0 {
				break
			}
			rest = rest[i+len(prefix):]
			end := strings.IndexAny(rest, `"'`)
			if end < 0 {
				t.Fatalf("unclosed asset ref after %q", prefix)
			}
			asset := "assets/" + rest[:end]
			if _, err := fs.Stat(sub, asset); err != nil {
				t.Fatalf("index references missing embedded asset %q: %v", asset, err)
			}
			rest = rest[end:]
		}
	}
}

func TestDeployMux_EmbeddedJSHasNoLocalhostWS(t *testing.T) {
	sub, err := fs.Sub(embeddedDist, "dist")
	if err != nil {
		t.Fatalf("embed subtree: %v", err)
	}
	entries, err := fs.ReadDir(sub, "assets")
	if err != nil {
		t.Fatalf("embedded assets missing — run `make build-client` before testing the deploy path: %v", err)
	}
	var jsCount int
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".js") {
			continue
		}
		jsCount++
		data, err := fs.ReadFile(sub, "assets/"+e.Name())
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		for _, bad := range []string{
			"localhost:8080",
			"ws://localhost",
			"wss://localhost",
			"VITE_WS_BASE_URL",
		} {
			if bytes.Contains(data, []byte(bad)) {
				t.Fatalf("embedded %s contains %q — production build must use same-origin WS", e.Name(), bad)
			}
		}
	}
	if jsCount == 0 {
		t.Fatal("no embedded JS assets — run `make build-client`")
	}
}

func TestDeployMux_QuestionContentStaysOffSpectatorWire(t *testing.T) {
	provider := deployQuestionProvider{question: room.Question{
		ID:     "q-deploy-secret",
		Type:   "mcq",
		Prompt: "DEPLOY_SECRET_PROMPT",
		Options: []room.QuestionOption{
			{ID: "DEPLOY_SECRET_CORRECT", Text: "DEPLOY_SECRET_CORRECT_TEXT", Correct: true},
			{ID: "DEPLOY_SECRET_WRONG", Text: "DEPLOY_SECRET_WRONG_TEXT"},
		},
	}}
	_, server := setupDeployServer(t, provider)

	connA := dialDeployWS(t, server)
	defer connA.Close()
	connB := dialDeployWS(t, server)
	defer connB.Close()

	deploySendJoin(t, connA, "Alice", "deploy-redact")
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	deploySendJoin(t, connB, "Bob", "deploy-redact")
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	_, _ = deployReadUntil(t, connB, ws.MessageTypeStateSync)

	choose, _ := json.Marshal(ws.ChooseLevelPayload{Difficulty: "easy"})
	if err := connA.WriteJSON(ws.Message{Type: ws.MessageTypeChooseLevel, RoomID: "deploy-redact", Payload: choose}); err != nil {
		t.Fatalf("choose_level: %v", err)
	}

	_, activeQ := deployReadUntil(t, connA, ws.MessageTypeQuestionStarted)
	if !bytes.Contains(activeQ.Payload, []byte("DEPLOY_SECRET_PROMPT")) {
		t.Fatalf("active player missing prompt: %s", activeQ.Payload)
	}
	spectatorRaw, spectatorQ := deployReadUntil(t, connB, ws.MessageTypeQuestionStarted)
	for _, secret := range []string{
		"DEPLOY_SECRET_PROMPT",
		"DEPLOY_SECRET_CORRECT",
		"DEPLOY_SECRET_CORRECT_TEXT",
		"DEPLOY_SECRET_WRONG",
		"DEPLOY_SECRET_WRONG_TEXT",
	} {
		if bytes.Contains(spectatorRaw, []byte(secret)) {
			t.Fatalf("deploy mux leaked %q to spectator: %s", secret, spectatorRaw)
		}
	}
	var spectatorPayload map[string]any
	if err := json.Unmarshal(spectatorQ.Payload, &spectatorPayload); err != nil {
		t.Fatalf("spectator payload: %v", err)
	}
	for _, forbidden := range []string{"problem_id", "type", "prompt", "options"} {
		if _, ok := spectatorPayload[forbidden]; ok {
			t.Fatalf("spectator payload has %q: %s", forbidden, spectatorQ.Payload)
		}
	}

	answer, _ := json.Marshal([]string{"DEPLOY_SECRET_CORRECT"})
	submit, _ := json.Marshal(ws.SubmitAnswerPayload{ProblemID: "q-deploy-secret", Answer: answer})
	if err := connA.WriteJSON(ws.Message{Type: ws.MessageTypeSubmitAnswer, RoomID: "deploy-redact", Payload: submit}); err != nil {
		t.Fatalf("submit: %v", err)
	}

	resultRaw, resultMsg := deployReadUntil(t, connB, ws.MessageTypeAnswerResult)
	for _, secret := range []string{"correct_answer", "DEPLOY_SECRET_CORRECT", "DEPLOY_SECRET_PROMPT"} {
		if bytes.Contains(resultRaw, []byte(secret)) {
			t.Fatalf("deploy mux answer_result leaked %q: %s", secret, resultRaw)
		}
	}
	var publicResult map[string]any
	if err := json.Unmarshal(resultMsg.Payload, &publicResult); err != nil {
		t.Fatalf("public result: %v", err)
	}
	if _, ok := publicResult["correct_answer"]; ok {
		t.Fatalf("spectator saw correct_answer: %s", resultMsg.Payload)
	}
}

func TestDeployMux_TimeoutThenLateSubmitResolvesOnce(t *testing.T) {
	provider := deployQuestionProvider{question: room.Question{
		ID:     "q-deploy-timeout",
		Type:   "mcq",
		Prompt: "timeout path",
		Options: []room.QuestionOption{
			{ID: "correct", Text: "yes", Correct: true},
			{ID: "wrong", Text: "no"},
		},
	}}
	hub, server := setupDeployServer(t, provider)

	connA := dialDeployWS(t, server)
	defer connA.Close()
	connB := dialDeployWS(t, server)
	defer connB.Close()

	const roomID = "deploy-timeout"
	deploySendJoin(t, connA, "Alice", roomID)
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	deploySendJoin(t, connB, "Bob", roomID)
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	_, _ = deployReadUntil(t, connB, ws.MessageTypeStateSync)

	r := hub.GetRoomInstance(roomID)
	r.SetDeadlineDurations(40*time.Millisecond, 40*time.Millisecond, 40*time.Millisecond)

	choose, _ := json.Marshal(ws.ChooseLevelPayload{Difficulty: "easy"})
	if err := connA.WriteJSON(ws.Message{Type: ws.MessageTypeChooseLevel, RoomID: roomID, Payload: choose}); err != nil {
		t.Fatalf("choose_level: %v", err)
	}
	_, _ = deployReadUntil(t, connA, ws.MessageTypeQuestionStarted)

	raw, result := deployReadUntil(t, connB, ws.MessageTypeAnswerResult)
	var payload map[string]any
	if err := json.Unmarshal(result.Payload, &payload); err != nil {
		t.Fatalf("answer_result: %v", err)
	}
	timedOut, _ := payload["timed_out"].(bool)
	correct, _ := payload["correct"].(bool)
	if !timedOut || correct {
		t.Fatalf("expected timeout win, got %s", raw)
	}

	answer, _ := json.Marshal([]string{"correct"})
	submit, _ := json.Marshal(ws.SubmitAnswerPayload{ProblemID: "q-deploy-timeout", Answer: answer})
	_ = connA.WriteJSON(ws.Message{Type: ws.MessageTypeSubmitAnswer, RoomID: roomID, Payload: submit})

	// Late submit must not produce a second public answer_result.
	extra := 0
	_ = connB.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	for {
		_, msgRaw, err := connB.ReadMessage()
		if err != nil {
			break
		}
		var msg ws.Message
		if err := json.Unmarshal(msgRaw, &msg); err != nil {
			continue
		}
		if msg.Type == ws.MessageTypeAnswerResult {
			extra++
		}
	}
	if extra != 0 {
		t.Fatalf("late submit after timeout produced %d extra answer_result(s)", extra)
	}

	aliceID := ""
	for _, p := range hub.GetRoomPlayers(roomID) {
		if p.Name == "Alice" {
			aliceID = p.ID
		}
	}
	if r.GetActivePlayerID() == aliceID {
		t.Fatalf("timeout did not advance turn, still %s", aliceID)
	}
}

func TestDeployMux_SubmitWinsAndStaleTimerDoesNotDoubleResolve(t *testing.T) {
	provider := deployQuestionProvider{question: room.Question{
		ID:     "q-deploy-submit",
		Type:   "mcq",
		Prompt: "submit path",
		Options: []room.QuestionOption{
			{ID: "correct", Text: "yes", Correct: true},
			{ID: "wrong", Text: "no"},
		},
	}}
	hub, server := setupDeployServer(t, provider)

	connA := dialDeployWS(t, server)
	defer connA.Close()
	connB := dialDeployWS(t, server)
	defer connB.Close()

	const roomID = "deploy-submit"
	deploySendJoin(t, connA, "Alice", roomID)
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	deploySendJoin(t, connB, "Bob", roomID)
	_, _ = deployReadUntil(t, connA, ws.MessageTypeStateSync)
	_, _ = deployReadUntil(t, connB, ws.MessageTypeStateSync)

	r := hub.GetRoomInstance(roomID)
	r.SetDeadlineDurations(120*time.Millisecond, 120*time.Millisecond, 120*time.Millisecond)

	choose, _ := json.Marshal(ws.ChooseLevelPayload{Difficulty: "easy"})
	if err := connA.WriteJSON(ws.Message{Type: ws.MessageTypeChooseLevel, RoomID: roomID, Payload: choose}); err != nil {
		t.Fatalf("choose_level: %v", err)
	}
	_, _ = deployReadUntil(t, connA, ws.MessageTypeQuestionStarted)
	_, _ = deployReadUntil(t, connB, ws.MessageTypeQuestionStarted)

	answer, _ := json.Marshal([]string{"correct"})
	submit, _ := json.Marshal(ws.SubmitAnswerPayload{ProblemID: "q-deploy-submit", Answer: answer})
	if err := connA.WriteJSON(ws.Message{Type: ws.MessageTypeSubmitAnswer, RoomID: roomID, Payload: submit}); err != nil {
		t.Fatalf("submit: %v", err)
	}

	_, result := deployReadUntil(t, connB, ws.MessageTypeAnswerResult)
	var payload map[string]any
	if err := json.Unmarshal(result.Payload, &payload); err != nil {
		t.Fatalf("answer_result: %v", err)
	}
	correct, _ := payload["correct"].(bool)
	timedOut, _ := payload["timed_out"].(bool)
	if !correct || timedOut {
		t.Fatalf("expected submit win, got %+v", payload)
	}

	// Wait past the original deadline; the stopped timer must not fire again.
	time.Sleep(200 * time.Millisecond)
	extra := 0
	_ = connB.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	for {
		_, msgRaw, err := connB.ReadMessage()
		if err != nil {
			break
		}
		var msg ws.Message
		if err := json.Unmarshal(msgRaw, &msg); err != nil {
			continue
		}
		if msg.Type == ws.MessageTypeAnswerResult {
			extra++
		}
	}
	if extra != 0 {
		t.Fatalf("stale deadline produced %d extra answer_result(s)", extra)
	}

	aliceID := ""
	for _, p := range hub.GetRoomPlayers(roomID) {
		if p.Name == "Alice" {
			aliceID = p.ID
		}
	}
	if r.GetActivePlayerID() == aliceID {
		t.Fatalf("submit did not advance turn")
	}
}
