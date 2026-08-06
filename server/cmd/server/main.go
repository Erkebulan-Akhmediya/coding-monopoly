package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	"server/internal/admin"
	"server/internal/ws"
)

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/monopoly?sslmode=disable"
	}
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Create database pool: %v", err)
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Connect to database: %v", err)
	}

	adminHandler, err := admin.NewHandler(db, admin.ConfigFromEnv())
	if err != nil {
		log.Fatalf("Configure admin API: %v", err)
	}

	hub := ws.NewHub(ws.NewDBQuestionProvider(db))
	go hub.Run()

	http.HandleFunc("/ws", ws.Handler(hub))

	// Admin spectator WebSocket: token validated before upgrade, admin clients
	// cannot trigger choose_level or submit_answer even if they try.
	http.HandleFunc("/ws/admin", ws.AdminHandler(hub, adminHandler.ValidateToken))

	http.Handle("/admin", adminHandler)
	http.Handle("/admin/", adminHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Monopoly Server Running"))
	})

	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
