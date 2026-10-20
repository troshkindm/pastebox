// Package testutil содержит вспомогательные функции для интеграционных тестов.
package testutil

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool подключается к тестовой PostgreSQL (переменная TEST_DATABASE_URL), заново создаёт
// таблицу snippets по файлу миграции и возвращает пул. Без переменной тест пропускается.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set: integration test skipped")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(ctx, `DROP TABLE IF EXISTS snippets`); err != nil {
		t.Fatalf("drop table: %v", err)
	}
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

// CountSnippets возвращает число записей в таблице snippets.
func CountSnippets(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM snippets`).Scan(&n); err != nil {
		t.Fatalf("count snippets: %v", err)
	}
	return n
}
