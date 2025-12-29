# Panduan Membuat Service Baru (Contoh: Customer)

Dokumen ini menjelaskan langkah-langkah untuk membuat service baru dalam project ini, dengan menggunakan studi kasus service `customer`. Arsitektur yang digunakan adalah **Hexagonal/Clean Architecture** dengan pemisahan **Domain**, **Usecase**, dan **Adapters**.

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
