package repository

import (
	"context"
	"errors"

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

// Create сохраняет текст и возвращает UUID новой записи.
func (r *SnippetRepository) Create(ctx context.Context, content string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO snippets (content) VALUES ($1) RETURNING id::text`, content).Scan(&id)
	return id, err
}

// Get возвращает запись по UUID или ErrNotFound.
func (r *SnippetRepository) Get(ctx context.Context, id string) (model.Snippet, error) {
	var s model.Snippet
	err := r.pool.QueryRow(ctx,
		`SELECT id::text, content, created_at FROM snippets WHERE id = $1::uuid`, id).
		Scan(&s.ID, &s.Content, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Snippet{}, ErrNotFound
	}
	return s, err
}
