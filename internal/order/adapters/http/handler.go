package http

import (
	"encoding/json"
	"net/http"

	"hex-postgres-grpc/internal/auth"
	order "hex-postgres-grpc/internal/order/domain"
)

type Handler struct {
	svc  order.Service
	auth auth.Service
}

func NewHandler(svc order.Service, authSvc auth.Service) *Handler {
	return &Handler{
		svc:  svc,
		auth: authSvc,
	}
}

type CreateOrderRequest struct {
	Amount float64 `json:"amount"`
}

type UpdateOrderRequest struct {
	Amount float64 `json:"amount"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("GET /orders/{id}", h.GetOrder)
	mux.HandleFunc("PUT /orders/{id}", h.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", h.DeleteOrder)
	mux.HandleFunc("GET /orders", h.ListOrders)
}

// CreateOrder creates a new order
// @Summary Create Order
// @Description Create a new order with the given amount
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Create Order Request"
// @Success 200 {object} order.Order
// @Failure 401 {string} string "unauthorized"
// @Router /orders [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := h.svc.CreateOrder(r.Context(), req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// GetOrder returns a single order by ID
// @Summary Get Order
// @Description Get details of a single order by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} order.Order
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	o, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		if err == order.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// UpdateOrder updates an existing order
// @Summary Update Order
// @Description Update the amount of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param request body UpdateOrderRequest true "Update Order Request"
// @Success 200 {object} order.Order
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "not found"
// @Router /orders/{id} [put]
func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := h.svc.UpdateOrder(r.Context(), id, req.Amount)
	if err != nil {
		if err == order.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// DeleteOrder deletes an order by ID
// @Summary Delete Order
// @Description Delete an order by ID
// @Tags orders
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "unauthorized"
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteOrder(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListOrders returns all orders
// @Summary List Orders
// @Description Get a list of all orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} order.Order
// @Failure 401 {string} string "unauthorized"
// @Router /orders [get]
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
