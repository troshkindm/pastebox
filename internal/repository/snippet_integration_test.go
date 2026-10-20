package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/troshkindm/pastebox/internal/repository"
	"github.com/troshkindm/pastebox/internal/testutil"
)

func TestCreateThenGet(t *testing.T) {
	repo := repository.NewSnippetRepository(testutil.Pool(t))
	ctx := context.Background()

	id, err := repo.Create(ctx, "integration text", nil)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := repo.Get(ctx, id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Content != "integration text" {
		t.Errorf("content: got %q, want %q", got.Content, "integration text")
	}
	if got.ID != id {
		t.Errorf("id: got %q, want %q", got.ID, id)
	}
	if got.CreatedAt.IsZero() {
		t.Errorf("created_at is not set")
	}
}

func TestGetMissing(t *testing.T) {
	repo := repository.NewSnippetRepository(testutil.Pool(t))

	_, err := repo.Get(context.Background(), "00000000-0000-4000-8000-000000000999")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error: got %v, want ErrNotFound", err)
	}
}

func TestGetExpiredAndAlive(t *testing.T) {
	repo := repository.NewSnippetRepository(testutil.Pool(t))
	ctx := context.Background()

	past := time.Now().Add(-time.Minute)
	expiredID, err := repo.Create(ctx, "old", &past)
	if err != nil {
		t.Fatalf("Create expired: %v", err)
	}
	if _, err := repo.Get(ctx, expiredID); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expired record: error got %v, want ErrNotFound", err)
	}

	future := time.Now().Add(time.Hour)
	aliveID, err := repo.Create(ctx, "fresh", &future)
	if err != nil {
		t.Fatalf("Create alive: %v", err)
	}
	got, err := repo.Get(ctx, aliveID)
	if err != nil {
		t.Fatalf("Get alive: %v", err)
	}
	if got.ExpiresAt == nil {
		t.Errorf("expires_at is not returned")
	}
}
