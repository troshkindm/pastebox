package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
