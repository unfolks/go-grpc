package http

import (
	"encoding/json"
	"net/http"

	product "hex-postgres-grpc/internal/product/domain"
	"hex-postgres-grpc/internal/product/usecase"
)

type Handler struct {
	createProduct *usecase.CreateProduct
	getProduct    *usecase.GetProduct
	listProducts  *usecase.ListProducts
	updateProduct *usecase.UpdateProduct
	deleteProduct *usecase.DeleteProduct
}

func NewHandler(
	createProduct *usecase.CreateProduct,
	getProduct *usecase.GetProduct,
	listProducts *usecase.ListProducts,
	updateProduct *usecase.UpdateProduct,
	deleteProduct *usecase.DeleteProduct,
) *Handler {
	return &Handler{
		createProduct: createProduct,
		getProduct:    getProduct,
		listProducts:  listProducts,
		updateProduct: updateProduct,
		deleteProduct: deleteProduct,
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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.createProduct.Execute(r.Context(), req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	p, err := h.getProduct.Execute(r.Context(), id)
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

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.updateProduct.Execute(r.Context(), id, req.Name, req.Price)
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

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	if err := h.deleteProduct.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.listProducts.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
