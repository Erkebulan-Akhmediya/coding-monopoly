package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

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

	hub := ws.NewHub(ws.NewDBQuestionProvider(db))
	go hub.Run()

	adminHandler, err := admin.NewHandler(db, admin.ConfigFromEnv(), hub)
	if err != nil {
		log.Fatalf("Configure admin API: %v", err)
	}

	handler := buildDeployMux(
		hub,
		adminHandler,
		ws.AdminHandler(hub, adminHandler.ValidateToken),
	)

	addr := listenAddr()
	log.Printf("Server listening on http://%s (LAN clients: use this host's IP)", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// listenAddr returns the TCP bind address.
// Prefer LISTEN_ADDR (e.g. 0.0.0.0:8080 or 10.10.40.69:8080); otherwise
// PORT (default 8080) bound on all interfaces so LAN clients can connect.
func listenAddr() string {
	if addr := strings.TrimSpace(os.Getenv("LISTEN_ADDR")); addr != "" {
		return addr
	}
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	return "0.0.0.0:" + port
}
