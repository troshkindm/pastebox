// Package testutil содержит вспомогательные функции для интеграционных тестов.
package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var schemaCounter atomic.Int64

// Pool подключается к тестовой PostgreSQL (переменная TEST_DATABASE_URL), создаёт для теста
// отдельную схему, применяет в ней миграцию и возвращает пул, работающий в этой схеме.
// Схемы изолируют тесты друг от друга: пакеты тестируются параллельно в одной базе.
// Без переменной тест пропускается.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set: integration test skipped")
	}
	ctx := context.Background()

	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	schema := fmt.Sprintf("t_%d_%d_%d", os.Getpid(), time.Now().UnixNano(), schemaCounter.Add(1))
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		admin.Close()
		t.Fatalf("create schema: %v", err)
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		admin.Close()
		t.Fatalf("parse test database url: %v", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		admin.Close()
		t.Fatalf("connect to test schema: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
		_, _ = admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
		admin.Close()
	})

	_, thisFile, _, _ := runtime.Caller(0)
	migration := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations", "001_create_snippets.sql")
	sql, err := os.ReadFile(migration)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}
	return pool
}

// CountSnippets возвращает число записей в таблице snippets тестовой схемы.
func CountSnippets(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM snippets`).Scan(&n); err != nil {
		t.Fatalf("count snippets: %v", err)
	}
	return n
}