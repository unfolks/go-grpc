package http

import (
	"encoding/json"
	"net/http"

	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/category/domain"
)

type Handler struct {
	service domain.Service
	auth    auth.Service
}

func NewHandler(service domain.Service, authSvc auth.Service) *Handler {
	return &Handler{
		service: service,
		auth:    authSvc,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /categories", h.CreateCategory)
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	mux.HandleFunc("PUT /categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)
	mux.HandleFunc("GET /categories", h.ListCategories)
}

// CreateCategory creates a new category
// @Summary Create Category
// @Description Create a new category with the provided name
// @Tags category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body struct{Name string `json:"name"`} true "Create Category Request"
// @Success 200 {object} domain.Category
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionCreate, auth.Resource{Type: "category"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := h.service.CreateCategory(r.Context(), req.Name, sub.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// GetCategory returns a single category by ID
// @Summary Get Category
// @Description Get a single category's details by their ID
// @Tags category
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} domain.Category
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "not found"
// @Router /categories/{id} [get]
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "category", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	cat, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// UpdateCategory updates an existing category
// @Summary Update Category
// @Description Update the name of an existing category
// @Tags category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body struct{Name string `json:"name"`} true "Update Category Request"
// @Success 200 {object} domain.Category
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "not found"
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionUpdate, auth.Resource{Type: "category", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cat, err := h.service.UpdateCategory(r.Context(), id, req.Name, sub.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// DeleteCategory deletes a category by ID
// @Summary Delete Category
// @Description Delete a category by ID
// @Tags category
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionDelete, auth.Resource{Type: "category", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = h.service.DeleteCategory(r.Context(), id, sub.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListCategories returns all categories
// @Summary List Categories
// @Description Get a list of all categories
// @Tags category
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Category
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "category"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	cats, err := h.service.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}
