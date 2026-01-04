package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hex-postgres-grpc/internal/auth"
	product "hex-postgres-grpc/internal/product/domain"
)

type Handler struct {
	service product.Service
	auth    auth.Service
}

func NewHandler(service product.Service, authSvc auth.Service) *Handler {
	return &Handler{
		service: service,
		auth:    authSvc,
	}
}

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.CreateProduct)
	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.HandleFunc("PUT /products/{id}", h.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProduct)
	mux.HandleFunc("GET /products", h.ListProducts)
}

// CreateProduct creates a new product
// @Summary Create Product
// @Description Create a new product with name and price
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Create Product Request"
// @Success 200 {object} product.Product
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionCreate, auth.Resource{Type: "product"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.CreateProduct(r.Context(), req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// GetProduct returns a single product by ID
// @Summary Get Product
// @Description Get details of a single product by ID
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} product.Product
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "not found"
// @Router /products/{id} [get]
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "product", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	p, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if err == product.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateProduct updates an existing product
// @Summary Update Product
// @Description Update name and price of an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body UpdateProductRequest true "Update Product Request"
// @Success 200 {object} product.Product
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "not found"
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionUpdate, auth.Resource{Type: "product", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.UpdateProduct(r.Context(), id, req.Name, req.Price)
	if err != nil {
		if err == product.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeleteProduct deletes a product by ID
// @Summary Delete Product
// @Description Delete a product by ID
// @Tags products
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionDelete, auth.Resource{Type: "product", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListProducts returns all products with pagination
// @Summary List Products
// @Description Get a list of all products with pagination
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} product.PaginatedResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /products [get]
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "product"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	resp, err := h.service.ListProductsPaginated(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
