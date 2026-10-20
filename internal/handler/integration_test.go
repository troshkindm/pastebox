package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/troshkindm/pastebox/internal/repository"
	"github.com/troshkindm/pastebox/internal/testutil"
)

// Сквозной сценарий handler -> repository -> PostgreSQL.
func TestIntegrationCreateAndView(t *testing.T) {
	pool := testutil.Pool(t)
	h := newHandler(t, repository.NewSnippetRepository(pool))

	rec := postPaste(h, "text from integration test")
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST /paste status: got %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if n := testutil.CountSnippets(t, pool); n != 1 {
		t.Fatalf("records in database: got %d, want 1", n)
	}

	view := do(h, httptest.NewRequest(http.MethodGet, rec.Header().Get("Location"), nil))
	if view.Code != http.StatusOK {
		t.Fatalf("GET status: got %d, want %d", view.Code, http.StatusOK)
	}
	if !strings.Contains(view.Body.String(), "text from integration test") {
		t.Errorf("page does not contain the saved text")
	}
}

func TestIntegrationEmptyInputWritesNothing(t *testing.T) {
	pool := testutil.Pool(t)
	h := newHandler(t, repository.NewSnippetRepository(pool))

	rec := postPaste(h, "   ")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if n := testutil.CountSnippets(t, pool); n != 0 {
		t.Errorf("records in database: got %d, want 0", n)
	}
}

func TestIntegrationMissingPaste(t *testing.T) {
	pool := testutil.Pool(t)
	h := newHandler(t, repository.NewSnippetRepository(pool))

	rec := do(h, httptest.NewRequest(http.MethodGet, "/p/00000000-0000-4000-8000-000000000999", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestIntegrationExpiredPasteIsNotFound(t *testing.T) {
	pool := testutil.Pool(t)
	repo := repository.NewSnippetRepository(pool)
	h := newHandler(t, repo)
	ctx := context.Background()

	past := time.Now().Add(-time.Hour)
	expired, err := repo.Create(ctx, "expired text", &past)
	if err != nil {
		t.Fatalf("Create expired: %v", err)
	}
	future := time.Now().Add(time.Hour)
	alive, err := repo.Create(ctx, "alive text", &future)
	if err != nil {
		t.Fatalf("Create alive: %v", err)
	}

	if rec := do(h, httptest.NewRequest(http.MethodGet, "/p/"+expired, nil)); rec.Code != http.StatusNotFound {
		t.Errorf("expired record: status got %d, want %d", rec.Code, http.StatusNotFound)
	}
	rec := do(h, httptest.NewRequest(http.MethodGet, "/p/"+alive, nil))
	if rec.Code != http.StatusOK {
		t.Errorf("alive record: status got %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "alive text") {
		t.Errorf("page does not contain the text of the alive record")
	}
}
