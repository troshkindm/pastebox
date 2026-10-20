package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/troshkindm/pastebox"
	"github.com/troshkindm/pastebox/internal/handler"
	"github.com/troshkindm/pastebox/internal/model"
	"github.com/troshkindm/pastebox/internal/repository"
)

// fakeStore — хранилище в памяти для модульных тестов обработчиков.
type fakeStore struct {
	snippets  map[string]model.Snippet
	createErr error
	creates   int
}

func newFakeStore() *fakeStore {
	return &fakeStore{snippets: map[string]model.Snippet{}}
}

func (f *fakeStore) Create(_ context.Context, content string) (string, error) {
	if f.createErr != nil {
		return "", f.createErr
	}
	f.creates++
	id := fmt.Sprintf("00000000-0000-4000-8000-%012d", f.creates)
	f.snippets[id] = model.Snippet{ID: id, Content: content, CreatedAt: time.Now()}
	return id, nil
}

func (f *fakeStore) Get(_ context.Context, id string) (model.Snippet, error) {
	s, ok := f.snippets[id]
	if !ok {
		return model.Snippet{}, repository.ErrNotFound
	}
	return s, nil
}

func newHandler(t *testing.T, store handler.SnippetStore) http.Handler {
	t.Helper()
	h, err := handler.New(store, pastebox.FS)
	if err != nil {
		t.Fatalf("handler.New: %v", err)
	}
	return h.Routes()
}

func do(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func postPaste(h http.Handler, content string) *httptest.ResponseRecorder {
	form := url.Values{"content": {content}}
	req := httptest.NewRequest(http.MethodPost, "/paste", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return do(h, req)
}

func TestIndexReturnsForm(t *testing.T) {
	h := newHandler(t, newFakeStore())
	rec := do(h, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `action="/paste"`) {
		t.Errorf("body does not contain the paste form")
	}
}

func TestCreatePasteValid(t *testing.T) {
	store := newFakeStore()
	h := newHandler(t, store)
	rec := postPaste(h, "hello world")

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/p/") {
		t.Errorf("Location: got %q, want prefix /p/", loc)
	}
	if store.creates != 1 {
		t.Errorf("records created: got %d, want 1", store.creates)
	}
}

func TestCreatePasteEmptyOrBlank(t *testing.T) {
	for _, content := range []string{"", "   ", "\n\t \n"} {
		store := newFakeStore()
		h := newHandler(t, store)
		rec := postPaste(h, content)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("content %q: status got %d, want %d", content, rec.Code, http.StatusBadRequest)
		}
		if store.creates != 0 {
			t.Errorf("content %q: record was created, want none", content)
		}
	}
}

func TestCreatePasteStoreError(t *testing.T) {
	store := newFakeStore()
	store.createErr = fmt.Errorf("database is down")
	h := newHandler(t, store)
	rec := postPaste(h, "hello")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestShowPasteExisting(t *testing.T) {
	store := newFakeStore()
	id, _ := store.Create(context.Background(), "saved text <b>")
	h := newHandler(t, store)
	rec := do(h, httptest.NewRequest(http.MethodGet, "/p/"+id, nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "saved text &lt;b&gt;") {
		t.Errorf("body does not contain the escaped text, got: %s", body)
	}
	if !strings.Contains(body, "14 симв.") {
		t.Errorf("body does not show the text length")
	}
}

func TestShowPasteMissing(t *testing.T) {
	h := newHandler(t, newFakeStore())
	rec := do(h, httptest.NewRequest(http.MethodGet, "/p/00000000-0000-4000-8000-000000000999", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestShowPasteInvalidID(t *testing.T) {
	h := newHandler(t, newFakeStore())
	rec := do(h, httptest.NewRequest(http.MethodGet, "/p/not-a-uuid", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}
