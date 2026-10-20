package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/troshkindm/pastebox/internal/model"
)

// ErrNotFound возвращается, если записи с таким id нет.
var ErrNotFound = errors.New("snippet not found")

// SnippetRepository хранит фрагменты в PostgreSQL.
type SnippetRepository struct {
	pool *pgxpool.Pool
}

func NewSnippetRepository(pool *pgxpool.Pool) *SnippetRepository {
	return &SnippetRepository{pool: pool}
}

// Create сохраняет текст (с необязательным сроком хранения) и возвращает UUID новой записи.
func (r *SnippetRepository) Create(ctx context.Context, content string, expiresAt *time.Time) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO snippets (content, expires_at) VALUES ($1, $2) RETURNING id::text`, content, expiresAt).Scan(&id)
	return id, err
}

// Get возвращает запись по UUID или ErrNotFound (в том числе для записи с истёкшим сроком хранения).
func (r *SnippetRepository) Get(ctx context.Context, id string) (model.Snippet, error) {
	var s model.Snippet
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, content, created_at, expires_at FROM snippets
		 WHERE id = $1::uuid AND (expires_at IS NULL OR expires_at > now())`, id).
		Scan(&s.ID, &s.Content, &s.CreatedAt, &s.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Snippet{}, ErrNotFound
	}
	return s, err
}
