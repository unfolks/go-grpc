# Panduan Membuat Service Baru (Contoh: Customer)

Dokumen ini menjelaskan langkah-langkah untuk membuat service baru dalam project ini, dengan menggunakan studi kasus service `customer`. Arsitektur yang digunakan adalah **Hexagonal/Clean Architecture** dengan pemisahan **Domain**, **Usecase**, dan **Adapters**.

## 1. Persiapan Struktur Folder
Buat folder baru di dalam `internal/` dengan nama service yang diinginkan, misalnya `internal/customer`.
Di dalamnya, buat struktur sub-folder sebagai berikut:

```text
internal/customer/
├── domain/       # Menyimpan Entity dan Interface (Core Logic)
├── usecase/      # Menyimpan Business Logic per fitur (Create, Get, List, dll)
└── adapters/     # Menyimpan implementasi infrastruktur
    ├── postgres/ # Implementasi Repository ke Database
    └── http/     # Implementasi HTTP Handler
```

---

## 2. Layer Domain (`internal/customer/domain`)

Layer ini adalah pusat dari service. Tidak boleh bergantung pada layer lain (no external imports selain stdlib atau uuid).

### a. Membuat Entity (`entity.go`)
Definisikan struct data utama.

```go
package domain

import "time"

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
```

### b. Membuat Interface Repository (`repository.go`)
Definisikan interface untuk komunikasi ke database. Implementation detail ada di layer adapter nanti.

```go
package domain

import "context"

type Repository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	// Tambahkan method lain sesuai kebutuhan
}
```

### c. Error Definition (Optional)
Jika ada error spesifik domain, definisikan di sini.

```go
var ErrNotFound = errors.New("customer not found")
```

---

## 3. Layer Usecase (`internal/customer/usecase`)

Layer ini berisi logika aplikasi. Setiap aksi (fitur) sebaiknya dibuat dalam file terpisah (Single Responsibility Principle).

### Contoh: Create Customer (`create_customer.go`)

```go
package usecase

import (
	"context"
	"time"
	"github.com/google/uuid"
	
	"hex-postgres-grpc/internal/customer/domain" // Import domain package
)

type CreateCustomer struct {
	repo domain.Repository
}

func NewCreateCustomer(repo domain.Repository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (u *CreateCustomer) Execute(ctx context.Context, name string, email string) (domain.Customer, error) {
    // Business Logic / Validation
	if name == "" {
		return domain.Customer{}, errors.New("name is required")
	}

	id := uuid.NewString()
	customer := domain.Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Save(ctx, &customer); err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}
```

*Lakukan hal yang sama untuk fitur lain seperti `get_customer.go`, `list_customers.go`, dll.*

---

## 4. Layer Adapters (`internal/customer/adapters`)

Layer ini menghubungkan aplikasi dengan dunia luar (Database, HTTP, gRPC).

### a. Postgres Repository (`adapters/postgres/repository.go`)
Implementasikan interface `domain.Repository`.

```go
package postgres

import (
	"context"
	"database/sql"
	"hex-postgres-grpc/internal/customer/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, c *domain.Customer) error {
	query := "INSERT INTO customers (id, name, email, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.Email, c.CreatedAt)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
    // Implementasi query select...
    return nil, nil // Placeholder
}
```

### b. HTTP Handler (`adapters/http/handler.go`)
Menangani request HTTP dan memanggil Usecase.

```go
package http

import (
	"encoding/json"
	"net/http"
	"hex-postgres-grpc/internal/customer/usecase"
)

type Handler struct {
	createCustomer *usecase.CreateCustomer
    // tambahkan usecase lain
}

func NewHandler(createCustomer *usecase.CreateCustomer) *Handler {
	return &Handler{createCustomer: createCustomer}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // Parsing Request
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Call Usecase
    customer, err := h.createCustomer.Execute(r.Context(), req.Name, req.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Response
    json.NewEncoder(w).Encode(customer)
}
```

---

## 5. Wiring (Setup Dependency Injection)

Terakhir, sambungkan semua komponen di `cmd/main.go` (atau file entry point aplikasi).

```go
// Di dalam fungsi main()

// 1. Init DB connection
db, _ := sql.Open("postgres", dbConnString)

// 2. Init Repository
customerRepo := customerPostgres.NewRepository(db)

// 3. Init Usecases
createCustomerUC := customerUsecase.NewCreateCustomer(customerRepo)

// 4. Init Handlers
customerHandler := customerHttp.NewHandler(createCustomerUC)

// 5. Register Routes
mux := http.NewServeMux()
customerHandler.RegisterRoutes(mux)

// 6. Start Server
http.ListenAndServe(":8080", mux)
```


---

## 6. Adapter gRPC (Optional)

Jika ingin menambahkan endpoint gRPC.

### a. Buat Definisi Proto
Buat folder `proto/customer/` dan file `customer.proto` di root project (sejajar `internal`).

```protobuf
syntax = "proto3";

package customerpb;
option go_package = "hex-postgres-grpc/proto/customer;customerpb";

import "google/protobuf/timestamp.proto";

service CustomerService {
    rpc CreateCustomer (CreateCustomerRequest) returns (CreateCustomerResponse);
}

message CustomerMessage {
    string id = 1;
    string name = 2;
    string email = 3;
    google.protobuf.Timestamp created_at = 4;
}

message CreateCustomerRequest {
    string name = 1;
    string email = 2;
}

message CreateCustomerResponse {
    CustomerMessage customer = 1;
}
```

### b. Generate Code

Pastikan tool `protoc-gen-go` dan `protoc-gen-go-grpc` sudah terinstall:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Pastikan folder binary Go ada di PATH (biasanya `$(go env GOPATH)/bin` atau `~/go/bin`).

Jalankan command berikut dari root project:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go_grpc_out=. --go_grpc_opt=paths=source_relative \
    proto/customer/customer.proto
```

### c. Implementasi Server gRPC (`internal/customer/adapters/grpc/server.go`)

```go
package grpc

import (
	"context"
	"hex-postgres-grpc/internal/customer/usecase"
	customerpb "hex-postgres-grpc/proto/customer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	createCustomer *usecase.CreateCustomer
}

func NewServer(createCustomer *usecase.CreateCustomer) *Server {
	return &Server{createCustomer: createCustomer}
}

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	c, err := s.createCustomer.Execute(ctx, req.Name, req.Email)
	if err != nil {
		return nil, err
	}
	
	return &customerpb.CreateCustomerResponse{
		Customer: &customerpb.CustomerMessage{
			Id:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			CreatedAt: timestamppb.New(c.CreatedAt),
		},
	}, nil
}
```

### d. Wiring gRPC Server (`cmd/main.go`)

```go
// ... di dalam wiring ...

// 1. Init gRPC Adapter
customerGRPC := customergrpc.NewServer(createCustomerUC)

// 2. Init gRPC Server Listener
lis, _ := net.Listen("tcp", ":50051")
grpcServer := grpc.NewServer()

// 3. Register Service
customerpb.RegisterCustomerServiceServer(grpcServer, customerGRPC)

// 4. Start Server
go func() {
    grpcServer.Serve(lis)
}()
```
