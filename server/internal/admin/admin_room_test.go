package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/ws"
)

type stubRoomManager struct {
	rooms map[string]bool
}

func (s *stubRoomManager) GetRoomsSummary() []ws.RoomSummary {
	out := make([]ws.RoomSummary, 0, len(s.rooms))
	for id := range s.rooms {
		out = append(out, ws.RoomSummary{RoomID: id})
	}
	return out
}

func (s *stubRoomManager) CreateRoom(roomID string) error {
	if s.rooms[roomID] {
		return ws.ErrRoomAlreadyExists
	}
	s.rooms[roomID] = true
	return nil
}

func TestCreateRoomEndpoint(t *testing.T) {
	now := time.Now()
	manager := &stubRoomManager{rooms: make(map[string]bool)}
	h := &Handler{
		config:      Config{Password: "secret", TokenSecret: "secret", TokenTTL: time.Minute},
		roomCreator: manager,
		roomLister:  manager,
		now:         func() time.Time { return now },
	}

	token, err := h.signToken(now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	body, _ := json.Marshal(map[string]string{"room_id": "lab-1"})
	req := httptest.NewRequest(http.MethodPost, "/admin/rooms", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create room status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !manager.rooms["lab-1"] {
		t.Fatal("room was not registered with manager")
	}

	dup := httptest.NewRequest(http.MethodPost, "/admin/rooms", bytes.NewReader(body))
	dup.Header.Set("Authorization", "Bearer "+token)
	dupRec := httptest.NewRecorder()
	h.ServeHTTP(dupRec, dup)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate create status=%d body=%s", dupRec.Code, dupRec.Body.String())
	}
}
