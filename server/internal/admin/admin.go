// Package admin implements the authenticated content-management HTTP API.
package admin

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTokenTTL = 15 * time.Minute

type Config struct {
	Password    string
	TokenSecret string
	TokenTTL    time.Duration
}

func ConfigFromEnv() Config {
	ttl := defaultTokenTTL
	if raw := os.Getenv("ADMIN_TOKEN_TTL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			ttl = parsed
		}
	}
	return Config{Password: os.Getenv("ADMIN_PASSWORD"), TokenSecret: os.Getenv("ADMIN_TOKEN_SECRET"), TokenTTL: ttl}
}

type Handler struct {
	db     *pgxpool.Pool
	config Config
	now    func() time.Time
}

func NewHandler(db *pgxpool.Pool, config Config) (*Handler, error) {
	if db == nil {
		return nil, errors.New("admin database pool is required")
	}
	if config.Password == "" {
		return nil, errors.New("ADMIN_PASSWORD must be configured")
	}
	if config.TokenSecret == "" {
		return nil, errors.New("ADMIN_TOKEN_SECRET must be configured")
	}
	if config.TokenTTL <= 0 {
		config.TokenTTL = defaultTokenTTL
	}
	return &Handler{db: db, config: config, now: time.Now}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/login" {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		h.login(w, r)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/admin/") || !h.authorized(r) {
		writeError(w, http.StatusUnauthorized, "admin authentication required")
		return
	}
	h.problems(w, r)
}

type loginRequest struct {
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if subtle.ConstantTimeCompare([]byte(request.Password), []byte(h.config.Password)) != 1 {
		writeError(w, http.StatusUnauthorized, "invalid admin password")
		return
	}
	expiresAt := h.now().Add(h.config.TokenTTL)
	token, err := h.signToken(expiresAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "issue admin token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "expires_at": expiresAt.UTC()})
}

func (h *Handler) signToken(expiresAt time.Time) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	payload, err := json.Marshal(map[string]any{"exp": expiresAt.Unix(), "nonce": base64.RawURLEncoding.EncodeToString(nonce)})
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(h.config.TokenSecret))
	_, _ = mac.Write([]byte(encoded))
	return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (h *Handler) authorized(r *http.Request) bool {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}
	pieces := strings.Split(parts[1], ".")
	if len(pieces) != 2 {
		return false
	}
	provided, err := base64.RawURLEncoding.DecodeString(pieces[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.config.TokenSecret))
	_, _ = mac.Write([]byte(pieces[0]))
	if !hmac.Equal(provided, mac.Sum(nil)) {
		return false
	}
	payload, err := base64.RawURLEncoding.DecodeString(pieces[0])
	if err != nil {
		return false
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if json.Unmarshal(payload, &claims) != nil || claims.Exp <= h.now().Unix() {
		return false
	}
	return true
}

type Option struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}
type Problem struct {
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	Difficulty      string    `json:"difficulty"`
	Title           string    `json:"title"`
	Prompt          string    `json:"prompt"`
	IsPublished     bool      `json:"is_published"`
	Options         []Option  `json:"options,omitempty"`
	AcceptedAnswers []string  `json:"accepted_answers,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type problemInput struct {
	Type        string `json:"type"`
	Difficulty  string `json:"difficulty"`
	Title       string `json:"title"`
	Prompt      string `json:"prompt"`
	IsPublished bool   `json:"is_published"`
	Options     []struct {
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"options"`
	AcceptedAnswers []string `json:"accepted_answers"`
}

func (h *Handler) problems(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/admin/problems")
	if path == "" || path == "/" {
		if r.Method == http.MethodGet {
			h.listProblems(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.createProblem(w, r)
			return
		}
		methodNotAllowed(w)
		return
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 2 && parts[1] == "publish" {
		if r.Method == http.MethodPost {
			h.publishProblem(w, r, parts[0])
			return
		}
		methodNotAllowed(w)
		return
	}
	if len(parts) != 1 {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	switch r.Method {
	case http.MethodPut:
		h.updateProblem(w, r, parts[0])
	case http.MethodDelete:
		h.deleteProblem(w, r, parts[0])
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) listProblems(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	typeFilter, err := enumFilter(q.Get("type"), "type", []string{"mcq", "text"})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	difficulty, err := enumFilter(q.Get("difficulty"), "difficulty", []string{"easy", "medium", "hard"})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var published any
	if raw, ok := q["is_published"]; ok {
		if len(raw) != 1 {
			writeError(w, http.StatusBadRequest, "is_published must be a boolean")
			return
		}
		value, err := strconv.ParseBool(raw[0])
		if err != nil {
			writeError(w, http.StatusBadRequest, "is_published must be a boolean")
			return
		}
		published = value
	}
	rows, err := h.db.Query(r.Context(), `SELECT id, type, difficulty, title, prompt, is_published, created_at, updated_at
		FROM problems WHERE ($1::problem_type IS NULL OR type = $1::problem_type)
		AND ($2::difficulty_level IS NULL OR difficulty = $2::difficulty_level)
		AND ($3::boolean IS NULL OR is_published = $3) ORDER BY created_at DESC`, typeFilter, difficulty, published)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list problems")
		return
	}
	defer rows.Close()
	problems := make([]Problem, 0)
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.ID, &p.Type, &p.Difficulty, &p.Title, &p.Prompt, &p.IsPublished, &p.CreatedAt, &p.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "read problems")
			return
		}
		problems = append(problems, p)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "read problems")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"problems": problems})
}

func enumFilter(value, name string, allowed []string) (any, error) {
	if value == "" {
		return nil, nil
	}
	for _, candidate := range allowed {
		if value == candidate {
			return value, nil
		}
	}
	return nil, fmt.Errorf("invalid %s", name)
}

func (h *Handler) createProblem(w http.ResponseWriter, r *http.Request) {
	var input problemInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := input.validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, 500, "create problem")
		return
	}
	defer tx.Rollback(ctx)
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO problems (type, difficulty, title, prompt, is_published) VALUES ($1, $2, $3, $4, $5) RETURNING id`, input.Type, input.Difficulty, input.Title, input.Prompt, input.IsPublished).Scan(&id)
	if err == nil {
		err = replaceContents(ctx, tx, id, input)
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeError(w, 500, "create problem")
		return
	}
	p, err := h.loadProblem(ctx, id)
	if err != nil {
		writeError(w, 500, "read created problem")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) updateProblem(w http.ResponseWriter, r *http.Request, id string) {
	var input problemInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := input.validate(); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	ctx := r.Context()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, 500, "update problem")
		return
	}
	defer tx.Rollback(ctx)
	// Delete dependent rows before changing type; the database deliberately
	// enforces that option/answer rows always match their problem's type.
	if _, err = tx.Exec(ctx, `DELETE FROM problem_options WHERE problem_id = $1`, id); err == nil {
		_, err = tx.Exec(ctx, `DELETE FROM problem_accepted_answers WHERE problem_id = $1`, id)
	}
	if err != nil {
		writeError(w, 500, "update problem")
		return
	}
	result, err := tx.Exec(ctx, `UPDATE problems SET type = $1, difficulty = $2, title = $3, prompt = $4, is_published = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6`, input.Type, input.Difficulty, input.Title, input.Prompt, input.IsPublished, id)
	if err == nil && result.RowsAffected() == 0 {
		writeError(w, 404, "problem not found")
		return
	}
	if err == nil {
		err = replaceContents(ctx, tx, id, input)
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeError(w, 500, "update problem")
		return
	}
	p, err := h.loadProblem(ctx, id)
	if err != nil {
		writeError(w, 500, "read updated problem")
		return
	}
	writeJSON(w, 200, p)
}

func replaceContents(ctx context.Context, tx pgx.Tx, id string, input problemInput) error {
	if input.Type == "mcq" {
		for _, option := range input.Options {
			if _, err := tx.Exec(ctx, `INSERT INTO problem_options (problem_id, text, is_correct) VALUES ($1, $2, $3)`, id, option.Text, option.IsCorrect); err != nil {
				return err
			}
		}
	} else {
		for _, answer := range input.AcceptedAnswers {
			if _, err := tx.Exec(ctx, `INSERT INTO problem_accepted_answers (problem_id, answer_text) VALUES ($1, $2)`, id, answer); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Handler) deleteProblem(w http.ResponseWriter, r *http.Request, id string) {
	result, err := h.db.Exec(r.Context(), `DELETE FROM problems WHERE id = $1`, id)
	if err != nil {
		writeError(w, 500, "delete problem")
		return
	}
	if result.RowsAffected() == 0 {
		writeError(w, 404, "problem not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) publishProblem(w http.ResponseWriter, r *http.Request, id string) {
	problem, err := h.loadProblem(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "problem not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read problem")
		return
	}
	if err := problem.validateForPublish(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.db.Exec(r.Context(), `UPDATE problems SET is_published = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		writeError(w, 500, "publish problem")
		return
	}
	if result.RowsAffected() == 0 {
		writeError(w, 404, "problem not found")
		return
	}
	p, err := h.loadProblem(r.Context(), id)
	if err != nil {
		writeError(w, 500, "read published problem")
		return
	}
	writeJSON(w, 200, p)
}

func (p Problem) validateForPublish() error {
	input := problemInput{Type: p.Type, Difficulty: p.Difficulty, Title: p.Title, Prompt: p.Prompt, AcceptedAnswers: p.AcceptedAnswers}
	for _, option := range p.Options {
		input.Options = append(input.Options, struct {
			Text      string `json:"text"`
			IsCorrect bool   `json:"is_correct"`
		}{Text: option.Text, IsCorrect: option.IsCorrect})
	}
	return input.validate()
}

func (h *Handler) loadProblem(ctx context.Context, id string) (Problem, error) {
	var p Problem
	err := h.db.QueryRow(ctx, `SELECT id, type, difficulty, title, prompt, is_published, created_at, updated_at FROM problems WHERE id = $1`, id).Scan(&p.ID, &p.Type, &p.Difficulty, &p.Title, &p.Prompt, &p.IsPublished, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, err
	}
	if p.Type == "mcq" {
		rows, err := h.db.Query(ctx, `SELECT id, text, is_correct FROM problem_options WHERE problem_id = $1 ORDER BY id`, id)
		if err != nil {
			return p, err
		}
		defer rows.Close()
		for rows.Next() {
			var option Option
			if err := rows.Scan(&option.ID, &option.Text, &option.IsCorrect); err != nil {
				return p, err
			}
			p.Options = append(p.Options, option)
		}
		return p, rows.Err()
	}
	rows, err := h.db.Query(ctx, `SELECT answer_text FROM problem_accepted_answers WHERE problem_id = $1 ORDER BY id`, id)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var answer string
		if err := rows.Scan(&answer); err != nil {
			return p, err
		}
		p.AcceptedAnswers = append(p.AcceptedAnswers, answer)
	}
	return p, rows.Err()
}

func (input problemInput) validate() error {
	if input.Type != "mcq" && input.Type != "text" {
		return errors.New("type must be mcq or text")
	}
	if input.Difficulty != "easy" && input.Difficulty != "medium" && input.Difficulty != "hard" {
		return errors.New("difficulty must be easy, medium, or hard")
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Prompt) == "" {
		return errors.New("title and prompt are required")
	}
	if input.Type == "mcq" {
		if len(input.Options) < 2 {
			return errors.New("mcq requires at least two options")
		}
		correct := false
		for _, option := range input.Options {
			if strings.TrimSpace(option.Text) == "" {
				return errors.New("option text is required")
			}
			correct = correct || option.IsCorrect
		}
		if !correct {
			return errors.New("mcq requires at least one correct option")
		}
		return nil
	}
	if len(input.AcceptedAnswers) == 0 {
		return errors.New("text requires at least one accepted answer")
	}
	for _, answer := range input.AcceptedAnswers {
		if strings.TrimSpace(answer) == "" {
			return errors.New("accepted answer is required")
		}
	}
	return nil
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("invalid JSON body")
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("invalid JSON body")
	}
	return nil
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
