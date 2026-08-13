package main

import (
	"net/http"

	"server/internal/ws"
)

// buildDeployMux wires the same routes the shipped binary uses: WS, optional
// admin HTTP/WS, health, then embedded Vue SPA last so API routes win.
func buildDeployMux(hub *ws.Hub, adminHTTP http.Handler, adminWS http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", ws.Handler(hub))

	if adminWS != nil {
		mux.Handle("/ws/admin", adminWS)
	}
	if adminHTTP != nil {
		mux.Handle("/admin", adminHTTP)
		mux.Handle("/admin/", adminHTTP)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Vue SPA + static assets (embedded via embed.FS). Registered last so API
	// and WebSocket routes always win.
	mux.Handle("/", staticHandler())

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
