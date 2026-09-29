package handler

import (
	"context"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/troshkindm/pastebox/internal/model"
	"github.com/troshkindm/pastebox/internal/repository"
)

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// SnippetStore — то, что обработчикам нужно от хранилища.
type SnippetStore interface {
	Create(ctx context.Context, content string) (string, error)
	Get(ctx context.Context, id string) (model.Snippet, error)
}

type Handler struct {
	store SnippetStore
	index *template.Template
	paste *template.Template
	fsys  fs.FS
}

// New разбирает шаблоны из fsys (каталоги templates и static) и возвращает обработчик.
func New(store SnippetStore, fsys fs.FS) (*Handler, error) {
	index, err := template.ParseFS(fsys, "templates/index.html")
	if err != nil {
		return nil, err
	}
	paste, err := template.ParseFS(fsys, "templates/paste.html")
	if err != nil {
		return nil, err
	}
	return &Handler{store: store, index: index, paste: paste, fsys: fsys}, nil
}

// Routes регистрирует маршруты приложения.
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.showIndex)
	r.Post("/paste", h.createPaste)
	r.Get("/p/{id}", h.showPaste)
	r.Handle("/static/*", http.FileServerFS(h.fsys))
	return r
}

type indexData struct {
	Error string
}

func (h *Handler) showIndex(w http.ResponseWriter, _ *http.Request) {
	h.render(w, h.index, http.StatusOK, indexData{})
}

func (h *Handler) createPaste(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	content := r.PostFormValue("content")
	if strings.TrimSpace(content) == "" {
		h.render(w, h.index, http.StatusBadRequest, indexData{Error: "Текст не должен быть пустым"})
		return
	}
	id, err := h.store.Create(r.Context(), content)
	if err != nil {
		log.Printf("create paste: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/p/"+id, http.StatusSeeOther)
}

func (h *Handler) showPaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		http.NotFound(w, r)
		return
	}
	s, err := h.store.Get(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("get paste %s: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.render(w, h.paste, http.StatusOK, s)
}

func (h *Handler) render(w http.ResponseWriter, t *template.Template, status int, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := t.Execute(w, data); err != nil {
		log.Printf("render: %v", err)
	}
}
