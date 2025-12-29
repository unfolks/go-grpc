# Hexagonal Architecture with Go, gRPC, and PostgreSQL

This project demonstrates a Hexagonal Architecture (Ports and Adapters) implementation in Go, featuring gRPC and HTTP interfaces, and a PostgreSQL database.

## Project Overview

The project is structured using a Modular Monolith architecture, where each business capability is encapsulated in its own module.
- **Modules**: Each feature has its own folder (e.g., `internal/order`, `internal/product`, `internal/customer`) containing its domain logic, adapters, and its own wiring.
- **App**: Orchestrates the modules (`internal/app`).

Each module typically contains:
- `domain/`: Business logic, entities, and interfaces.
- `adapters/`: Implementations for databases, HTTP, gRPC, etc.
- `app.go`: Module-level wiring (initializes domain services and adapters).

## Prerequisites

- **Go**: 1.25+
- **PostgreSQL**: Database instance.
- **Protoc**: Protocol Buffer compiler (for regenerating proto files).
- **Go Plugins**:
    - `protoc-gen-go`
    - `protoc-gen-go-grpc`

## Setup & Installation

1.  **Clone the repository**
2.  **Install Dependencies**
    ```bash
    go mod download
    ```
3.  **Database Setup**
    Ensure you have a PostgreSQL database running. Update the configuration in `cmd/server/main.go`.
    
    *Default Config in `main.go`:*
    - User: `admin`
    - Password: `password123`
    - Host: `localhost:5432`
    - DB Name: `ordersdb`

    *Database Schema:*
    ```sql
    CREATE TABLE orders (
        id VARCHAR(36) PRIMARY KEY,
        amount DOUBLE PRECISION NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

    CREATE TABLE products (
        id VARCHAR(36) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        price DOUBLE PRECISION NOT NULL
    );

    CREATE TABLE customers (
        id VARCHAR(36) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        address TEXT
    );
    ```

## Usage

### Running with Makefile

The project includes a `Makefile` for convenience.

- **Run the server**: `make run`
- **Build the binary**: `make build` (outputs to `bin/server`)
- **Run tests**: `make test`
- **Tidy dependencies**: `make tidy`
- **Clean build artifacts**: `make clean`

### Running Manually
```bash
go run cmd/server/main.go
```
The server starts:
- HTTP Server on port `:8080`
- gRPC Server on port `:50051`

### API Endpoints

#### HTTP
- **Orders**
    - `POST /orders`: Create an order
    - `GET /orders/get?id=<ID>`: Get an order
- **Products**
    - `POST /products`: Create a product
    - `GET /products/{id}`: Get a product
    - `PUT /products/{id}`: Update a product
    - `DELETE /products/{id}`: Delete a product
    - `GET /products`: List all products
- **Customers**
    - `POST /customer`: Create a customer
    - `GET /customer`: List all customers

#### gRPC
- **OrderService** (CreateOrder, GetOrder)
- **ProductService** (CreateProduct, GetProduct, ListProducts, UpdateProduct, DeleteProduct)
- **CustomerService** (CreateCustomer, ListCustomers)
- Proto definitions: `proto/*.proto`

## Development Guide: How to Create a New Module

Follow these steps to add a new module (e.g., `Payment`).

### 1. Create Module Structure
```
internal/payment/
├── domain/        (Entity, Repository Interface, Service Interface)
├── usecase/       (Service Implementation)
├── adapters/      (Postgres, HTTP, gRPC implementations)
└── app.go         (Module-level wiring)
```

### 2. Implement Module Wiring (`internal/payment/app.go`)
Define a `Components` struct and an `Init` function that returns it:
```go
package payment

type Components struct {
    HTTPHandler *http.Handler
    GRPCServer  *grpc.Server
}

func Init(db *sql.DB) Components {
    repo := postgres.NewRepository(db)
    service := usecase.NewService(repo)
    
    httpHandler := http.NewHandler(service)
    grpcServer := grpc.NewServer(service)

    return Components{
        HTTPHandler: httpHandler,
        GRPCServer:  grpcServer,
    }
}
```

### 3. Register in App Wiring (`internal/app/wiring.go`)
Add the new module to the `Application` struct and call its `Init` in `app.Init`:
```go
type Application struct {
    DB      *sql.DB
    Payment payment.Components
    // ... other modules
}

func Init(cfg DBConfig) (*Application, error) {
    // ... db init
    return &Application{
        DB:      db,
        Payment: payment.Init(db),
    }, nil
}
```

### 4. Expose in Main (`cmd/server/main.go`)
Register routes or gRPC services:
```go
a.Payment.HTTPHandler.RegisterRoutes(mux)
```

## Regeneration of Proto Files
If you modify `.proto` files, run:

## Testing with Insomnia

### HTTP Requests
1.  **Create Order**
    - Method: `POST`
    - URL: `http://localhost:8080/orders`
    - Body (JSON):
        ```json
        {
            "amount": 100.50
        }
        ```
2.  **Get Order**
    - Method: `GET`
    - URL: `http://localhost:8080/orders/get`
    - Query Parameter: `id` = `<YOUR_ORDER_ID>` (e.g., from the Create Order response)

### gRPC Requests
1.  Create a new Request and select **gRPC** as the type.
2.  Click on **"Add Proto File"** and select the `proto/order.proto` file from this project.
    - *Note: You might need to add the project root to the "Import Paths" if it complains about imports.*
3.  Select the Service: `ORderService` (Note the casing).
4.  **CreateOrder**
    - Method: `CreateOrder`
    - Body:
        ```json
        {
            "amount": 250.00
        }
        ```
5.  **GetOrder**
    - Method: `GetOrder`
    - Body:
        ```json
        {
            "id": "<YOUR_ORDER_ID>"
        }
        ```
