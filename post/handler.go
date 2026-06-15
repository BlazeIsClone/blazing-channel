package post

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/blazeisclone/blazing-channel/pkg/api"
)

type PostService interface {
	Create(ctx context.Context, cmd CreatePostCommand) (*Post, error)
	GetAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id int) (*Post, error)
	Update(ctx context.Context, id int, cmd UpdatePostCommand) (*Post, error)
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	svc PostService
}

func NewHandler(svc PostService) *Handler {
	return &Handler{svc: svc}
}

func Routes(mux *http.ServeMux, svc PostService) {
	h := NewHandler(svc)

	mux.HandleFunc("GET "+api.Path("v1", "/posts"), h.Index)
	mux.HandleFunc("POST "+api.Path("v1", "/posts"), h.Store)
	mux.HandleFunc("GET "+api.Path("v1", "/posts/{id}"), h.Show)
	mux.HandleFunc("PUT "+api.Path("v1", "/posts/{id}"), h.Update)
	mux.HandleFunc("DELETE "+api.Path("v1", "/posts/{id}"), h.Destroy)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.GetAll(r.Context())

	if err != nil {
		jsonError(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []Post{}
	}

	jsonResponse(w, posts, http.StatusOK)
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var req postRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := req.validate(); errs.HasErrors() {
		jsonValidationErrors(w, errs)
		return
	}

	post, err := h.svc.Create(r.Context(), CreatePostCommand{
		Title: strings.TrimSpace(req.Title),
		Body:  strings.TrimSpace(req.Body),
	})

	if err != nil {
		jsonError(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, post, http.StatusCreated)
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	post, err := h.svc.FindByID(r.Context(), id)

	if errors.Is(err, ErrNotFound) {
		jsonError(w, "post not found", http.StatusNotFound)
		return
	}

	if err != nil {
		jsonError(w, "failed to fetch post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, post, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)

	if !ok {
		return
	}

	var req postRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := req.validate(); errs.HasErrors() {
		jsonValidationErrors(w, errs)
		return
	}

	post, err := h.svc.Update(r.Context(), id, UpdatePostCommand{
		Title: strings.TrimSpace(req.Title),
		Body:  strings.TrimSpace(req.Body),
	})

	if errors.Is(err, ErrNotFound) {
		jsonError(w, "post not found", http.StatusNotFound)
		return
	}

	if err != nil {
		jsonError(w, "failed to update post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, post, http.StatusOK)
}

func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)

	if !ok {
		return
	}

	err := h.svc.Delete(r.Context(), id)

	if errors.Is(err, ErrNotFound) {
		jsonError(w, "post not found", http.StatusNotFound)
		return
	}

	if err != nil {
		jsonError(w, "failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func jsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, map[string]string{"error": msg}, status)
}

func jsonValidationErrors(w http.ResponseWriter, errs ValidationErrors) {
	jsonResponse(w, map[string]ValidationErrors{"errors": errs}, http.StatusUnprocessableEntity)
}
