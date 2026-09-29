package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/troshkindm/pastebox"
	"github.com/troshkindm/pastebox/internal/handler"
	"github.com/troshkindm/pastebox/internal/repository"
)

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	// Пул подключается к БД лениво: приложение стартует, даже если PostgreSQL ещё недоступен.
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("database config: %v", err)
	}
	defer pool.Close()

	h, err := handler.New(repository.NewSnippetRepository(pool), pastebox.FS)
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("PasteBox listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
