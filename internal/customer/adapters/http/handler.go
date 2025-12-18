package http

import (
	"encoding/json"
	"hex-postgres-grpc/internal/customer/usecase"
	"net/http"
)

type Handler struct {
	createCustomer *usecase.CreateCustomer
	listCustomer   *usecase.ListCustomer
}

func NewHandler(createCustomer *usecase.CreateCustomer, listCustomer *usecase.ListCustomer) *Handler {
	return &Handler{
		createCustomer: createCustomer,
		listCustomer:   listCustomer,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customer", h.Create)
	mux.HandleFunc("GET /customer", h.List)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	customers, err := h.listCustomer.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cust, err := h.createCustomer.Execute(r.Context(), req.Name, req.Email, req.Address)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cust)
}
