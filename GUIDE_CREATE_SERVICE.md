# Panduan Membuat Service Baru

Dokumen ini menjelaskan langkah-langkah untuk membuat service baru dalam project ini menggunakan **Hexagonal/Clean Architecture** dengan pemisahan **Domain**, **Usecase**, dan **Adapters**.

## Tentang Arsitektur Ini

Arsitektur ini mengikuti prinsip **Dependency Inversion** di mana:
- **Domain Layer** (core) tidak bergantung pada layer manapun
- **Usecase Layer** bergantung pada Domain
- **Adapters Layer** (infrastruktur) bergantung pada Domain dan Usecase

### Layer Responsibilities

| Layer | Tanggung Jawab | Tidak Boleh Mengandung |
|-------|----------------|------------------------|
| **Domain** | Entity, Interface Repository, Interface Service, Error definitions | Database logic, HTTP logic, External dependencies |
| **Usecase** | Business logic implementation, Subject tracking, Validation | Database queries, HTTP handling, Framework-specific code |
| **Adapters** | HTTP handlers, gRPC servers, Database repositories, Authorization checks | Business logic |

### Standar Production dalam Project Ini

Semua service production di project ini mengimplementasikan:
- ✅ **Audit Tracking** - created_by, updated_by, deleted_by
- ✅ **Soft Delete** - deleted_at untuk menandai penghapusan
- ✅ **Authorization** - ABAC (Attribute-Based Access Control)
- ✅ **Pagination** - untuk semua list endpoints
- ✅ **Context Propagation** - untuk tracking user dan tracing
- ✅ **Domain Errors** - error definitions di layer domain

### Struktur Dokumen

Dokumen ini terbagi menjadi:
1. **Basic Service (Customer)** - Contoh minimal untuk memahami konsep dasar
2. **Production Service (Product)** - Implementasi lengkap dengan semua pattern production
3. **Integration & Best Practices** - Panduan wiring dan deployment

---

# PART 1: BASIC SERVICE (Customer Example)

> **📘 Catatan**: Bagian ini menampilkan contoh minimal untuk memahami struktur dasar. Untuk implementasi production-ready dengan audit tracking, authorization, dan pagination, lihat **PART 2: PRODUCTION SERVICE (Product Example)**.

## 1. Persiapan Struktur Folder
Buat folder baru di dalam `internal/` dengan nama service yang diinginkan, misalnya `internal/customer`.
Di dalamnya, buat struktur sub-folder sebagai berikut:

```text
internal/customer/
├── domain/       # Menyimpan Entity dan Service Interface (Core Logic)
├── usecase/      # Menyimpan Implementasi Service (Business Logic)
├── adapters/     # Menyimpan implementasi infrastruktur
│   ├── postgres/ # Implementasi Repository ke Database
│   ├── http/     # Implementasi HTTP Handler
│   └── grpc/     # Implementasi gRPC Server
└── app.go        # Komponen wiring tingkat modul
```

---

## 2. Layer Domain (`internal/customer/domain`)

Layer ini adalah pusat dari service. Tidak boleh bergantung pada layer lain.

### a. Membuat Entity (`entity.go`)
Definisikan struct data utama.

```go
package domain

import "time"

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
```

### b. Membuat Interface Repository (`repository.go`)
Definisikan interface untuk komunikasi ke database.

```go
package domain

import "context"

type Repository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindAll(ctx context.Context) ([]Customer, error)
}
```

### c. Membuat Interface Service (`service.go`)
Definisikan kontrak business logic yang akan digunakan oleh adapter (HTTP/gRPC).

```go
package domain

import "context"

type Service interface {
	CreateCustomer(ctx context.Context, name, email, address string) (*Customer, error)
	ListCustomers(ctx context.Context) ([]Customer, error)
}
```

---

## 3. Layer Usecase (`internal/customer/usecase`)

Layer ini berisi implementasi dari `domain.Service`. Semua logika bisnis dikonsolidasikan dalam satu service implementation.

### Membuat Service Implementation (`service.go`)

```go
package usecase

import (
	"context"
	"time"
	"github.com/google/uuid"
	"hex-postgres-grpc/internal/customer/domain"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.Service {
	return &service{repo: repo}
}

func (s *service) CreateCustomer(ctx context.Context, name, email, address string) (*domain.Customer, error) {
	id := uuid.NewString()
	customer := domain.Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		Address:   address,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *service) ListCustomers(ctx context.Context) ([]domain.Customer, error) {
	return s.repo.FindAll(ctx)
}
```

---

## 4. Layer Adapters (`internal/customer/adapters`)

### a. HTTP Handler (`adapters/http/handler.go`)
Menangani request HTTP dan memanggil Service.

```go
package http

import (
	"encoding/json"
	"net/http"
	"hex-postgres-grpc/internal/customer/domain"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", h.Create)
	mux.HandleFunc("GET /customers", h.List)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name    string `json:"name"`
        Email   string `json:"email"`
        Address string `json:"address"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    customer, err := h.service.CreateCustomer(r.Context(), req.Name, req.Email, req.Address)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(customer)
}
```

### b. gRPC Server (`adapters/grpc/server.go`)
Implementasikan interface yang di-generate dari proto.

```go
package grpc

import (
	"context"
	"hex-postgres-grpc/internal/customer/domain"
	customerpb "hex-postgres-grpc/proto/customer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	service domain.Service
}

func NewServer(service domain.Service) *Server {
	return &Server{service: service}
}

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	c, err := s.service.CreateCustomer(ctx, req.Name, req.Email, req.Address)
	if err != nil {
		return nil, err
	}
	return &customerpb.CreateCustomerResponse{
		Customer: &customerpb.CustomerMessage{
			Id:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			Address:   c.Address,
			CreatedAt: timestamppb.New(c.CreatedAt),
		},
	}, nil
}
```

---

## 5. Module Wiring (`internal/customer/app.go`)

Inisialisasi semua komponen dalam satu tempat.

```go
package customer

import (
	"database/sql"
	"hex-postgres-grpc/internal/customer/adapters/grpc"
	"hex-postgres-grpc/internal/customer/adapters/http"
	"hex-postgres-grpc/internal/customer/adapters/postgres"
	"hex-postgres-grpc/internal/customer/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

func Init(db *sql.DB) Components {
	repo := postgres.NewRepository(db)
	service := usecase.NewService(repo)

	return Components{
		HTTPHandler: http.NewHandler(service),
		GRPCServer:  grpc.NewServer(service),
	}
}
```

---

## 6. Global Wiring & Registration

### a. `internal/app/wiring.go`
Tambahkan module baru ke struct `Application`.

```go
type Application struct {
	DB       *sql.DB
	Customer customer.Components
    // ...
}

func Init(cfg DBConfig) (*Application, error) {
    // ...
	return &Application{
		DB:       db,
		Customer: customer.Init(db),
	}, nil
}
```

### b. `cmd/server/main.go`
Daftarkan handler ke server HTTP dan gRPC.

```go
// HTTP
a.Customer.HTTPHandler.RegisterRoutes(mux)

// gRPC
customerpb.RegisterCustomerServiceServer(grpcServer, a.Customer.GRPCServer)
```

---

# PART 2: PRODUCTION SERVICE (Product Example)

> **🚀 Production-Ready**: Bagian ini menampilkan implementasi lengkap dengan semua pattern production yang digunakan dalam codebase ini, termasuk audit tracking, authorization, pagination, dan soft delete.

## 1. Struktur Folder Product Service

```text
internal/product/
├── domain/
│   ├── entity.go       # Entity dengan audit fields dan pagination models
│   ├── repository.go   # Repository interface dengan pagination
│   └── service.go      # Service interface
├── usecase/
│   └── service.go      # Service implementation dengan subject tracking
├── adapters/
│   ├── postgres/
│   │   └── repository.go  # PostgreSQL implementation
│   ├── http/
│   │   └── handler.go     # HTTP handler dengan authorization
│   └── grpc/
│       └── product_grpc.go # gRPC server
└── app.go              # Module wiring dengan auth service
```

---

## 2. Layer Domain dengan Production Patterns

### a. Entity dengan Audit Fields (`domain/entity.go`)

```go
package product

import (
	"errors"
	"time"
)

// Domain errors
var ErrNotFound = errors.New("product not found")
var ErrInvalidPrice = errors.New("invalid price")

// Product entity dengan audit fields lengkap
type Product struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`  // Nullable - hanya terisi saat update
	UpdatedBy *string    `json:"updated_by,omitempty"`  // Nullable
	DeletedAt *time.Time `json:"deleted_at,omitempty"`  // Soft delete marker
	DeletedBy *string    `json:"deleted_by,omitempty"`  // Soft delete
}

// PaginatedData untuk response pagination
type PaginatedData struct {
	Data             []Product `json:"data"`
	CurrentPage      int       `json:"current_page"`
	HaveNextPage     bool      `json:"have_next_page"`
	HavePreviousPage bool      `json:"have_previous_page"`
	Limit            int       `json:"limit"`
	TotalItem        int       `json:"total_item"`
	TotalPage        int       `json:"total_page"`
}

// PaginatedResponse wrapper untuk API response
type PaginatedResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    PaginatedData `json:"data"`
}
```

**Key Points:**
- **Audit Fields**: `created_by`, `updated_by`, `deleted_by` untuk tracking siapa yang melakukan operasi
- **Nullable Fields**: Gunakan pointer (`*time.Time`, `*string`) untuk field yang optional
- **Soft Delete**: `deleted_at` untuk menandai record yang dihapus tanpa menghilangkan data
- **Domain Errors**: Definisikan error di domain layer untuk business rules
- **Pagination Models**: Struct terpisah untuk response pagination yang konsisten

### b. Repository Interface dengan Pagination (`domain/repository.go`)

```go
package product

import "context"

type Repository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string, deletedBy string) error  // Soft delete dengan deletedBy
	FindAllPaginated(ctx context.Context, limit, offset int) ([]Product, int, error)
}
```

**Key Points:**
- Semua method menerima `context.Context` sebagai parameter pertama
- `Delete` method menerima `deletedBy` untuk audit tracking
- `FindAllPaginated` return total count untuk kalkulasi pagination
- Interface ini tidak tahu implementasi database (dependency inversion)

### c. Service Interface (`domain/service.go`)

```go
package product

import (
	"context"
)

type Service interface {
	CreateProduct(ctx context.Context, name string, price float64) (Product, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, id string, name string, price float64) (Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProductsPaginated(ctx context.Context, page, limit int) (PaginatedResponse, error)
}
```

**Key Points:**
- Context propagation untuk subject tracking dan tracing
- Return `PaginatedResponse` untuk list endpoints
- Interface mendefinisikan kontrak business logic
- Tidak ada implementation detail di interface

---

## 3. Layer Usecase dengan Subject Tracking

### Implementation dengan Audit Tracking (`usecase/service.go`)

```go
package usecase

import (
	"context"
	"time"

	"hex-postgres-grpc/internal/auth"
	product "hex-postgres-grpc/internal/product/domain"

	"github.com/google/uuid"
)

// SystemUserID digunakan sebagai fallback ketika tidak ada authenticated user
const SystemUserID = "00000000-0000-0000-0000-000000000000"

type service struct {
	repo product.Repository
}

func NewService(repo product.Repository) product.Service {
	return &service{repo: repo}
}

// CreateProduct dengan business validation dan audit tracking
func (s *service) CreateProduct(ctx context.Context, name string, price float64) (product.Product, error) {
	// Business validation
	if price < 0 {
		return product.Product{}, product.ErrInvalidPrice
	}

	// Extract subject from context untuk audit tracking
	sub, _ := auth.SubjectFromContext(ctx)
	createdBy := SystemUserID
	if sub.ID != "" {
		createdBy = sub.ID
	}

	// Create entity dengan audit fields
	id := uuid.NewString()
	p := product.Product{
		ID:        id,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,  // Track who created
	}

	if err := s.repo.Save(ctx, &p); err != nil {
		return product.Product{}, err
	}
	return p, nil
}

func (s *service) GetProduct(ctx context.Context, id string) (product.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}
	return *p, nil
}

// UpdateProduct dengan audit tracking
func (s *service) UpdateProduct(ctx context.Context, id string, name string, price float64) (product.Product, error) {
	// Fetch existing product
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}

	// Business validation
	if price < 0 {
		return product.Product{}, product.ErrInvalidPrice
	}

	// Extract subject untuk audit tracking
	sub, _ := auth.SubjectFromContext(ctx)
	updatedBy := SystemUserID
	if sub.ID != "" {
		updatedBy = sub.ID
	}

	// Update fields dengan audit tracking
	now := time.Now()
	p.Name = name
	p.Price = price
	p.UpdatedAt = &now       // Set update timestamp
	p.UpdatedBy = &updatedBy // Track who updated

	if err := s.repo.Update(ctx, p); err != nil {
		return product.Product{}, err
	}

	return *p, nil
}

// DeleteProduct (soft delete) dengan audit tracking
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	// Extract subject untuk audit tracking
	sub, _ := auth.SubjectFromContext(ctx)
	deletedBy := SystemUserID
	if sub.ID != "" {
		deletedBy = sub.ID
	}

	// Soft delete dengan tracking siapa yang menghapus
	return s.repo.Delete(ctx, id, deletedBy)
}

// ListProductsPaginated dengan pagination logic
func (s *service) ListProductsPaginated(ctx context.Context, page, limit int) (product.PaginatedResponse, error) {
	// Default values dan validation
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch data dengan pagination
	products, total, err := s.repo.FindAllPaginated(ctx, limit, offset)
	if err != nil {
		return product.PaginatedResponse{}, err
	}

	// Calculate total pages
	totalPage := 0
	if limit > 0 {
		totalPage = (total + limit - 1) / limit  // Ceil division
	}

	// Build response dengan pagination metadata
	return product.PaginatedResponse{
		Success: true,
		Message: "Get data product successfully",
		Data: product.PaginatedData{
			Data:             products,
			CurrentPage:      page,
			HaveNextPage:     page < totalPage,
			HavePreviousPage: page > 1,
			Limit:            limit,
			TotalItem:        total,
			TotalPage:        totalPage,
		},
	}, nil
}
```

**Key Patterns:**
1. **Subject Extraction**: `auth.SubjectFromContext(ctx)` untuk mendapatkan user yang melakukan request
2. **SystemUserID Fallback**: Digunakan ketika tidak ada authenticated user (system operations)
3. **Business Validation**: Validasi di usecase layer (price >= 0)
4. **Audit Tracking**: Populate `created_by`, `updated_by`, `deleted_by` dari subject
5. **Nullable Fields**: Gunakan pointer untuk `UpdatedAt` dan `UpdatedBy`
6. **Pagination Logic**: Calculate offset dan total pages

---

## 4. HTTP Handler dengan Authorization

### Handler dengan ABAC Authorization (`adapters/http/handler.go`)

```go
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
	auth    auth.Service  // Auth service dependency
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

// RegisterRoutes menggunakan Go 1.22+ routing pattern
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.CreateProduct)
	mux.HandleFunc("GET /products/{id}", h.GetProduct)
	mux.HandleFunc("PUT /products/{id}", h.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", h.DeleteProduct)
	mux.HandleFunc("GET /products", h.ListProducts)
}

// CreateProduct dengan authorization check
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
	// 1. Extract subject dari context (sudah di-set oleh auth middleware)
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Authorization check menggunakan ABAC
	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionCreate, auth.Resource{Type: "product"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 3. Parse request
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Call service layer
	p, err := h.service.CreateProduct(r.Context(), req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// GetProduct dengan authorization dan error mapping
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
	// 1. Extract path parameter (Go 1.22+)
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	// 2. Extract subject
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. Authorization check dengan resource ID
	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "product", ID: id})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 4. Call service
	p, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		// Error mapping: domain error → HTTP status code
		if err == product.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateProduct dengan authorization
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

	// Authorization dengan action UPDATE
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

// DeleteProduct dengan authorization
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

	// Authorization dengan action DELETE
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

// ListProducts dengan pagination dan authorization
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

	// Authorization untuk list (read action tanpa resource ID)
	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, auth.Resource{Type: "product"})
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Parse pagination query parameters
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
```

**Key Patterns:**
1. **Auth Service Injection**: Handler menerima `auth.Service` sebagai dependency
2. **Subject Extraction**: `auth.SubjectFromContext(r.Context())` di setiap handler
3. **Authorization Check**: `h.auth.Authorize(ctx, subject, action, resource)` sebelum business logic
4. **Action Constants**: `ActionCreate`, `ActionRead`, `ActionUpdate`, `ActionDelete`
5. **Resource Struct**: `auth.Resource{Type: "product", ID: id}` untuk ABAC
6. **Modern Routing**: `r.PathValue("id")` untuk path parameters (Go 1.22+)
7. **Error Mapping**: Domain errors (`ErrNotFound`) → HTTP status codes (404)
8. **Swagger Annotations**: `@Summary`, `@Security BearerAuth`, dll untuk OpenAPI documentation
9. **Pagination Params**: Extract `page` dan `limit` dari query parameters

---

## 5. PostgreSQL Repository dengan Soft Delete

### Repository Implementation (`adapters/postgres/repository.go`)

```go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	product "hex-postgres-grpc/internal/product/domain"
)

type ProductRepoPG struct {
	db *sql.DB
}

func NewProductRepoPG(db *sql.DB) *ProductRepoPG {
	return &ProductRepoPG{db: db}
}

// Save dengan audit fields
func (r *ProductRepoPG) Save(ctx context.Context, p *product.Product) error {
	const q = `INSERT INTO products (id, name, price, created_at, created_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Price, p.CreatedAt, p.CreatedBy)
	return err
}

// FindByID dengan soft delete filter
func (r *ProductRepoPG) FindByID(ctx context.Context, id string) (*product.Product, error) {
	// IMPORTANT: WHERE deleted_at IS NULL untuk filter soft deleted records
	const query = `SELECT id, name, price, created_at, created_by, updated_at, updated_by
	               FROM products
	               WHERE id = $1 AND deleted_at IS NULL`

	var p product.Product
	var created time.Time
	row := r.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&p.ID, &p.Name, &p.Price, &created, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy); err != nil {
		// Convert sql.ErrNoRows to domain error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, err
	}
	p.CreatedAt = created
	return &p, nil
}

// Update dengan audit fields dan soft delete filter
func (r *ProductRepoPG) Update(ctx context.Context, p *product.Product) error {
	const q = `UPDATE products
	           SET name = $1, price = $2, updated_at = $3, updated_by = $4
	           WHERE id = $5 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, p.Name, p.Price, p.UpdatedAt, p.UpdatedBy, p.ID)
	return err
}

// Delete (soft delete) dengan audit tracking
func (r *ProductRepoPG) Delete(ctx context.Context, id string, deletedBy string) error {
	const q = `UPDATE products
	           SET deleted_at = $1, deleted_by = $2
	           WHERE id = $3 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, time.Now(), deletedBy, id)
	return err
}

// FindAllPaginated dengan soft delete filter dan count
func (r *ProductRepoPG) FindAllPaginated(ctx context.Context, limit, offset int) ([]product.Product, int, error) {
	// Count total items (excluding soft deleted)
	const countQ = `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
	var total int
	if err := r.db.QueryRowContext(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated data
	const q = `SELECT id, name, price, created_at, created_by, updated_at, updated_by
	           FROM products
	           WHERE deleted_at IS NULL
	           LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []product.Product{}
	for rows.Next() {
		var p product.Product
		var created time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &created, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy); err != nil {
			return nil, 0, err
		}
		p.CreatedAt = created
		products = append(products, p)
	}
	return products, total, rows.Err()
}
```

**Key Patterns:**
1. **Soft Delete Filter**: Semua query menggunakan `WHERE deleted_at IS NULL`
2. **Error Conversion**: `sql.ErrNoRows` → `product.ErrNotFound` (domain error)
3. **NULL Handling**: Scan ke pointer untuk nullable fields (`UpdatedAt`, `UpdatedBy`)
4. **Pagination**: Separate COUNT query dan data query dengan LIMIT/OFFSET
5. **Context Propagation**: Gunakan `*Context` methods untuk cancellation dan timeout

---

## 6. Module Wiring dengan Auth Service

### Application Component (`app.go`)

```go
package product

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/product/adapters/grpc"
	"hex-postgres-grpc/internal/product/adapters/http"
	"hex-postgres-grpc/internal/product/adapters/postgres"
	"hex-postgres-grpc/internal/product/usecase"
)

type Components struct {
	HTTPHandler *http.Handler
	GRPCServer  *grpc.Server
}

// Init dengan dependency injection auth service
func Init(db *sql.DB, authSvc auth.Service) Components {
	// 1. Initialize repository
	repo := postgres.NewProductRepoPG(db)

	// 2. Initialize service (usecase)
	service := usecase.NewService(repo)

	// 3. Initialize adapters dengan dependencies
	httpHandler := http.NewHandler(service, authSvc)  // HTTP handler perlu auth service
	grpcServer := grpc.NewProductGRPCServer(service)   // gRPC menggunakan global interceptor

	return Components{
		HTTPHandler: httpHandler,
		GRPCServer:  grpcServer,
	}
}
```

**Key Points:**
- Auth service di-inject ke HTTP handler untuk authorization checks
- gRPC server tidak perlu auth service karena menggunakan global interceptor
- Dependency injection pattern untuk testability

---

# PART 3: INTEGRATION & DEPLOYMENT

## 1. Authorization Integration

### A. Auth Service Interface

Auth service menyediakan dua fungsi utama:
1. **Authentication**: Memvalidasi JWT token dan extract subject
2. **Authorization**: Memeriksa apakah subject boleh melakukan action pada resource

```go
package auth

import "context"

// Subject represents authenticated user
type Subject struct {
	ID   string
	Role string
}

// Resource untuk ABAC
type Resource struct {
	Type string  // e.g., "product", "user"
	ID   string  // optional, untuk specific resource
}

// Action constants
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

type Service interface {
	// Authorize memeriksa apakah subject dapat melakukan action pada resource
	Authorize(ctx context.Context, subject Subject, action string, resource Resource) (bool, error)
}
```

### B. Global Middleware Setup

Di `cmd/server/main.go`, auth middleware diterapkan secara global:

```go
func main() {
	// ... setup database, config, etc ...

	// Initialize application dengan semua modules
	app, err := application.Init(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer app.DB.Close()

	// === HTTP SERVER (port 8080) ===
	mux := http.NewServeMux()

	// Register all module routes
	app.AuthHandler.RegisterRoutes(mux)        // Auth routes (login, register)
	app.Product.HTTPHandler.RegisterRoutes(mux) // Product routes

	// IMPORTANT: Wrap dengan auth middleware secara global
	handler := app.Auth.HTTPMiddleware(mux)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,  // Semua request melalui auth middleware
	}

	// === gRPC SERVER (port 50051) ===
	// IMPORTANT: Register auth interceptor secara global
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(app.Auth.GRPCUnaryInterceptor),
	)

	// Register all gRPC services
	productpb.RegisterProductServiceServer(grpcServer, app.Product.GRPCServer)

	// ... start servers ...
}
```

**Key Points:**
- Auth middleware di-apply ke **semua routes** via `app.Auth.HTTPMiddleware(mux)`
- gRPC interceptor di-register saat create server dengan `grpc.UnaryInterceptor()`
- Subject sudah tersedia di context ketika sampai ke handler
- Handler hanya perlu extract subject dan check authorization

### C. Authorization Check Pattern

Setiap handler mengikuti pattern ini:

```go
func (h *Handler) SomeOperation(w http.ResponseWriter, r *http.Request) {
	// 1. Extract subject (sudah di-set oleh middleware)
	sub, ok := auth.SubjectFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Check authorization
	resource := auth.Resource{
		Type: "product",     // Resource type
		ID:   "some-id",     // Optional: specific resource ID
	}
	authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead, resource)
	if err != nil || !authorized {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 3. Proceed dengan business logic
	// ...
}
```

### D. ABAC Policy Structure

Authorization service menggunakan ABAC (Attribute-Based Access Control):

```go
// Example policy structure
type Policy struct {
	Role         string  // e.g., "admin", "user"
	Action       string  // "create", "read", "update", "delete"
	Resource     string  // "product", "user"
	AllowedIf    func(subject Subject, resource Resource) bool
}

// Example policies:
// - Admin dapat melakukan semua action pada semua resource
// - User hanya dapat read product
// - User dapat update/delete product miliknya sendiri (check resource.ID == subject.ID)
```

---

## 2. Application-Level Wiring

### A. Update `internal/app/wiring.go`

```go
package application

import (
	"database/sql"
	"hex-postgres-grpc/internal/auth"
	authhttp "hex-postgres-grpc/internal/auth/adapters/http"
	"hex-postgres-grpc/internal/product"
	// import modules lain...
)

type Application struct {
	DB          *sql.DB
	Auth        auth.Components    // Auth module (harus diinit pertama)
	AuthHandler *authhttp.Handler  // Auth HTTP handler
	Product     product.Components // Product module
	// modules lain...
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Init(cfg DBConfig) (*Application, error) {
	// 1. Setup database connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 2. IMPORTANT: Initialize auth module FIRST
	authComponents := auth.Init(db)

	// 3. Initialize other modules dengan auth service dependency
	productComponents := product.Init(db, authComponents.Service)

	return &Application{
		DB:          db,
		Auth:        authComponents,
		AuthHandler: authComponents.HTTPHandler,
		Product:     productComponents,
	}, nil
}
```

**Key Points:**
1. **Auth module diinit pertama** karena module lain depend on auth service
2. **Pass auth service** ke module lain via `Init(db, authSvc)`
3. **Application struct** menyimpan semua components untuk di-register di main.go

### B. Module Registration di `cmd/server/main.go`

```go
// HTTP Routes Registration
mux := http.NewServeMux()
app.AuthHandler.RegisterRoutes(mux)           // /login, /register
app.Product.HTTPHandler.RegisterRoutes(mux)   // /products, /products/{id}
// Register routes dari module lain...

// Wrap dengan auth middleware
handler := app.Auth.HTTPMiddleware(mux)

// gRPC Services Registration
grpcServer := grpc.NewServer(
	grpc.UnaryInterceptor(app.Auth.GRPCUnaryInterceptor),
)
productpb.RegisterProductServiceServer(grpcServer, app.Product.GRPCServer)
// Register gRPC services dari module lain...
```

---

## 3. Database Schema & Migrations

### A. Products Table Schema

```sql
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DOUBLE PRECISION NOT NULL,

    -- Audit fields
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by VARCHAR(36) NULL,

    -- Soft delete fields
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(36) NULL
);

-- Index untuk query performance
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_products_created_by ON products(created_by);
CREATE INDEX idx_products_name ON products(name);
```

**Key Points:**
1. **Audit columns**: `created_by`, `updated_by`, `deleted_by` untuk tracking
2. **Nullable columns**: `updated_at`, `updated_by`, `deleted_at`, `deleted_by` (NULL sampai digunakan)
3. **Index on deleted_at**: Untuk performance query `WHERE deleted_at IS NULL`
4. **ID format**: VARCHAR(36) untuk UUID

### B. Migration Best Practices

**Struktur Migration Files:**
```text
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_products_table.up.sql
├── 002_create_products_table.down.sql
└── ...
```

**Example Migration (`002_create_products_table.up.sql`):**
```sql
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DOUBLE PRECISION NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by VARCHAR(36) NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by VARCHAR(36) NULL
);

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
```

**Down Migration (`002_create_products_table.down.sql`):**
```sql
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP TABLE IF EXISTS products;
```

---

## 4. Best Practices & Patterns

### A. Kapan Menggunakan Audit Fields

**✅ GUNAKAN untuk:**
- Entity bisnis utama (products, users, orders, dll)
- Data yang perlu compliance/audit trail
- Data yang akan di-modify oleh multiple users
- Data yang penting untuk business operations

**❌ JANGAN GUNAKAN untuk:**
- Lookup tables/reference data yang jarang berubah
- Cache tables atau temporary data
- Tables yang hanya diakses system (bukan user)

### B. Soft Delete vs Hard Delete

**Gunakan SOFT DELETE untuk:**
- Data bisnis yang mungkin perlu di-restore
- Data yang referenced oleh data lain
- Compliance requirements (harus keep history)
- User-facing entities (products, users, orders)

**Gunakan HARD DELETE untuk:**
- Session data, temporary data
- Cache entries
- Log files yang sudah tidak relevan
- Data yang truly tidak diperlukan lagi

**Implementation Soft Delete:**
```sql
-- Query harus SELALU include WHERE deleted_at IS NULL
SELECT * FROM products WHERE deleted_at IS NULL;

-- Soft delete
UPDATE products SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2;

-- Restore (jika diperlukan)
UPDATE products SET deleted_at = NULL, deleted_by = NULL WHERE id = $1;
```

### C. Pagination Strategy

**Best Practices:**
1. **Selalu paginate** untuk list endpoints
2. **Default values**: page=1, limit=10
3. **Max limit**: Set maximum (e.g., 100) untuk prevent large queries
4. **Return metadata**: total_items, total_pages, has_next, has_previous
5. **Offset-based** untuk simple cases, **cursor-based** untuk large datasets

**Implementation:**
```go
func (s *service) ListPaginated(ctx context.Context, page, limit int) (PaginatedResponse, error) {
	// Validation
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100  // Max limit
	}

	offset := (page - 1) * limit

	// Fetch dengan separate COUNT query
	items, total, err := s.repo.FindAllPaginated(ctx, limit, offset)
	// ... build response
}
```

### D. Error Handling Hierarchy

**Layer-by-layer error handling:**

1. **Domain Layer**: Define business errors
   ```go
   var ErrNotFound = errors.New("product not found")
   var ErrInvalidPrice = errors.New("invalid price")
   ```

2. **Repository Layer**: Convert infrastructure errors to domain errors
   ```go
   if errors.Is(err, sql.ErrNoRows) {
       return nil, product.ErrNotFound
   }
   ```

3. **Usecase Layer**: Handle business logic errors
   ```go
   if price < 0 {
       return product.Product{}, product.ErrInvalidPrice
   }
   ```

4. **Handler Layer**: Map domain errors to HTTP status codes
   ```go
   if err == product.ErrNotFound {
       http.Error(w, err.Error(), http.StatusNotFound)
       return
   }
   if err == product.ErrInvalidPrice {
       http.Error(w, err.Error(), http.StatusBadRequest)
       return
   }
   http.Error(w, err.Error(), http.StatusInternalServerError)
   ```

### E. Nullable Fields Pattern

**Kapan menggunakan pointer:**
```go
type Product struct {
	ID        string     `json:"id"`              // Required: tidak nullable
	Name      string     `json:"name"`            // Required: tidak nullable
	CreatedAt time.Time  `json:"created_at"`      // Required: selalu ada
	UpdatedAt *time.Time `json:"updated_at,omitempty"`  // Optional: NULL sampai ada update
	UpdatedBy *string    `json:"updated_by,omitempty"`  // Optional: NULL sampai ada update
}
```

**Setting nullable fields:**
```go
now := time.Now()
updatedBy := "user-id"
product.UpdatedAt = &now
product.UpdatedBy = &updatedBy
```

### F. Context Propagation

**Best Practices:**
1. **Selalu pass context** sebagai parameter pertama
2. **Jangan store context** di struct
3. **Propagate sampai ke database layer** untuk timeout/cancellation
4. **Gunakan context untuk values** (subject, trace ID, dll)

```go
// ✅ GOOD
func (s *service) DoSomething(ctx context.Context, arg string) error {
	return s.repo.Save(ctx, data)
}

// ❌ BAD
func (s *service) DoSomething(arg string) error {
	return s.repo.Save(data)
}
```

### G. Subject Tracking Pattern

**Extract subject di usecase layer:**
```go
const SystemUserID = "00000000-0000-0000-0000-000000000000"

func (s *service) CreateProduct(ctx context.Context, name string, price float64) (Product, error) {
	// Extract subject dari context
	sub, _ := auth.SubjectFromContext(ctx)

	// Fallback ke SystemUserID jika tidak ada
	createdBy := SystemUserID
	if sub.ID != "" {
		createdBy = sub.ID
	}

	product := Product{
		// ...
		CreatedBy: createdBy,
	}
	// ...
}
```

**SystemUserID Use Cases:**
- Background jobs/cron tasks
- System-initiated operations
- Migration scripts
- Seed data

---

## 5. Quick Reference Checklist

### Checklist untuk Membuat Service Baru

#### 1. Setup Struktur Folder
- [ ] Buat folder `internal/{service_name}/`
- [ ] Buat subfolder: `domain/`, `usecase/`, `adapters/postgres/`, `adapters/http/`, `adapters/grpc/`
- [ ] Buat file `app.go` di root service folder

#### 2. Domain Layer
- [ ] **entity.go**: Define entity dengan audit fields (`created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`)
- [ ] **entity.go**: Define pagination models (`PaginatedData`, `PaginatedResponse`)
- [ ] **entity.go**: Define domain errors (`ErrNotFound`, `ErrInvalidX`)
- [ ] **repository.go**: Define repository interface dengan `ctx context.Context` di semua method
- [ ] **repository.go**: Include pagination method (`FindAllPaginated`)
- [ ] **service.go**: Define service interface dengan `ctx context.Context`

#### 3. Usecase Layer
- [ ] **service.go**: Implement service interface
- [ ] **service.go**: Add `SystemUserID` constant
- [ ] **service.go**: Extract subject dari context: `auth.SubjectFromContext(ctx)`
- [ ] **service.go**: Populate audit fields (`created_by`, `updated_by`, `deleted_by`)
- [ ] **service.go**: Implement business validations
- [ ] **service.go**: Implement pagination logic

#### 4. PostgreSQL Repository
- [ ] **repository.go**: Implement repository interface
- [ ] **repository.go**: Add soft delete filter di SEMUA queries: `WHERE deleted_at IS NULL`
- [ ] **repository.go**: Convert `sql.ErrNoRows` → domain error
- [ ] **repository.go**: Implement pagination dengan separate COUNT query
- [ ] **repository.go**: Handle NULL values dengan pointer scan

#### 5. HTTP Handler
- [ ] **handler.go**: Inject auth service dependency
- [ ] **handler.go**: Extract subject di setiap handler
- [ ] **handler.go**: Check authorization sebelum business logic
- [ ] **handler.go**: Use `r.PathValue()` untuk path parameters
- [ ] **handler.go**: Map domain errors → HTTP status codes
- [ ] **handler.go**: Add Swagger annotations
- [ ] **handler.go**: Parse pagination query params untuk list endpoints

#### 6. gRPC Handler (Optional)
- [ ] **{service}_grpc.go**: Implement gRPC service
- [ ] **{service}_grpc.go**: Extract subject dari context (global interceptor handles auth)

#### 7. Module Wiring
- [ ] **app.go**: Create `Components` struct
- [ ] **app.go**: Implement `Init(db, authSvc)` function
- [ ] **app.go**: Wire repository → service → handlers

#### 8. Application Integration
- [ ] **internal/app/wiring.go**: Add module ke `Application` struct
- [ ] **internal/app/wiring.go**: Initialize module di `Init()` dengan auth service
- [ ] **cmd/server/main.go**: Register HTTP routes
- [ ] **cmd/server/main.go**: Register gRPC service

#### 9. Database
- [ ] Create migration file dengan audit columns
- [ ] Add indexes: `deleted_at`, `created_by`
- [ ] Add constraints: CHECK, NOT NULL, FOREIGN KEY
- [ ] Run migration

#### 10. Testing (Recommended)
- [ ] Unit tests untuk usecase layer
- [ ] Integration tests untuk repository
- [ ] HTTP handler tests dengan mock auth
- [ ] Test soft delete behavior
- [ ] Test pagination edge cases

---

## 6. Common Patterns Quick Reference

### Pattern: Authorization Check
```go
sub, ok := auth.SubjectFromContext(r.Context())
if !ok {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
	return
}

authorized, err := h.auth.Authorize(r.Context(), sub, auth.ActionRead,
	auth.Resource{Type: "product", ID: id})
if err != nil || !authorized {
	http.Error(w, "forbidden", http.StatusForbidden)
	return
}
```

### Pattern: Subject Extraction
```go
sub, _ := auth.SubjectFromContext(ctx)
userID := SystemUserID
if sub.ID != "" {
	userID = sub.ID
}
```

### Pattern: Soft Delete Query
```go
const query = `SELECT * FROM products WHERE id = $1 AND deleted_at IS NULL`
```

### Pattern: Pagination Calculation
```go
offset := (page - 1) * limit
totalPage := (total + limit - 1) / limit  // Ceil division
hasNext := page < totalPage
```

### Pattern: Error Mapping
```go
if err == domain.ErrNotFound {
	return 404
} else if err == domain.ErrInvalidInput {
	return 400
} else {
	return 500
}
```

### Pattern: Nullable Field Assignment
```go
now := time.Now()
userID := "some-id"
entity.UpdatedAt = &now
entity.UpdatedBy = &userID
```

---

## 7. Troubleshooting

### Issue: "unauthorized" error di setiap request
**Solusi:**
- Pastikan auth middleware terpasang di main.go
- Cek JWT token format: `Authorization: Bearer <token>`
- Verify middleware extract subject dengan benar

### Issue: Authorization selalu forbidden
**Solusi:**
- Check policy configuration di auth service
- Verify subject.Role sesuai dengan policy
- Debug `Authorize()` method untuk lihat policy matching

### Issue: Soft deleted records masih muncul
**Solusi:**
- Pastikan SEMUA query include `WHERE deleted_at IS NULL`
- Check repository methods satu per satu
- Gunakan grep untuk find missing filters: `grep -r "FROM products" --include="*.go"`

### Issue: Pagination total count salah
**Solusi:**
- Pastikan COUNT query juga include `WHERE deleted_at IS NULL`
- Verify count query executed sebelum data query
- Check total calculation: `(total + limit - 1) / limit`

### Issue: Updated_at tidak ter-set
**Solusi:**
- Pastikan assign pointer bukan value: `entity.UpdatedAt = &now` bukan `entity.UpdatedAt = now`
- Check database column type (harus NULLABLE)
- Verify INSERT/UPDATE query include updated_at column

---

## Summary

### Basic vs Production Service

| Aspect | Basic (Customer) | Production (Product) |
|--------|------------------|----------------------|
| **Audit Fields** | ❌ Tidak ada | ✅ created_by, updated_by, deleted_by |
| **Soft Delete** | ❌ Hard delete | ✅ Soft delete dengan deleted_at |
| **Authorization** | ❌ Tidak ada | ✅ ABAC di setiap endpoint |
| **Pagination** | ❌ Return semua data | ✅ Paginated dengan metadata |
| **Subject Tracking** | ❌ Tidak ada | ✅ Extract dari context |
| **Domain Errors** | ❌ Generic errors | ✅ Specific business errors |
| **Error Mapping** | ❌ Tidak ada | ✅ Domain → HTTP status codes |
| **Context Propagation** | ⚠️ Partial | ✅ Full propagation |

### Kapan Menggunakan Pattern Mana?

**Gunakan Basic Pattern untuk:**
- Prototype/POC
- Internal tools yang tidak perlu audit
- Simple CRUD tanpa business rules kompleks

**Gunakan Production Pattern untuk:**
- Production applications
- Services dengan multiple users
- Compliance requirements
- Business-critical data
- Services yang perlu audit trail

---

**Selamat! 🎉** Anda sekarang memiliki panduan lengkap untuk membuat service baru dengan production-ready patterns yang digunakan di codebase ini.
