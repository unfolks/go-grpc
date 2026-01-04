package http

import (
	"encoding/json"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/customer/domain"
	"net/http"
)

type Handler struct {
	service domain.Service
	auth    auth.Service
}

func NewHandler(service domain.Service, auth auth.Service) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customer", h.Create)
	mux.HandleFunc("GET /customer", h.List)
	mux.HandleFunc("GET /customer/:id", h.Get)
}

// List returns all customers
// @Summary List Customers
// @Description Get a list of all customers
// @Tags customer
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Customer
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /customer [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "customer"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	customers, err := h.service.ListCustomers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// Create creates a new customer
// @Summary Create Customer
// @Description Create a new customer with the provided details
// @Tags customer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,email=string,address=string} true "Create Customer Request"
// @Success 200 {object} domain.Customer
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Router /customer [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionCreate, auth.Resource{Type: "customer"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cust, err := h.service.CreateCustomer(r.Context(), req.Name, req.Email, req.Address)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cust)
}

// Get returns a single customer by ID
// @Summary Get Customer
// @Description Get a single customer's details by their ID
// @Tags customer
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.Customer
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 404 {string} string "not found"
// @Router /customer/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	cust, err := h.service.GetCustomer(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ABAC check: Owner or Admin
	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{
		Type: "customer",
		ID:   id,
		Attributes: map[string]interface{}{
			"owner_id": cust.ID, // Assuming customer ID is the owner ID for this example
		},
	})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cust)
}
