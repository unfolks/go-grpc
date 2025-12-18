# Hexagonal Architecture with Go, gRPC, and PostgreSQL

This project demonstrates a Hexagonal Architecture (Ports and Adapters) implementation in Go, featuring gRPC and HTTP interfaces, and a PostgreSQL database.

## Project Overview

The project is structured using a Modular Monolith architecture, grouping code by feature (module) rather than technical layer.
- **Modules**: Each feature has its own folder (e.g., `internal/order`, `internal/product`) containing its domain logic and adapters.
- **App**: Wires everything together (`internal/app`).

Each module typically contains:
- `domain/`: Business logic, entities, and interfaces.
- `adapters/`: Implementations for databases, HTTP, gRPC, etc.

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
    Ensure you have a PostgreSQL database running. Update the configuration in `cmd/server/main.go` or use environment variables (recommended for production).
    
    *Default Config in `main.go`:*
    - User: `admin`
    - Password: `password123`
    - Host: `db:5432`
    - DB Name: `ordersdb`

    *Note: You will need to create the `orders` table manually or via migration:*
    ```sql
    CREATE TABLE orders (
        id VARCHAR(36) PRIMARY KEY,
        amount DOUBLE PRECISION NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
    ```

## Usage

### Running the Server
```bash
go run cmd/server/main.go
```
The server starts:
- HTTP Server on port `:8080`
- gRPC Server on port `:50051`

### API Endpoints

#### HTTP
- **Create Order**
    - `POST /orders`
    - Body: `{"amount": 100.50}`
- **Get Order**
    - `GET /orders/get?id=<ORDER_ID>`

#### gRPC
- **CreateOrder**
- **GetOrder**
- Proto definition: `proto/order.proto`

## Development Guide: How to Create a New Module

Follow these steps to add a new module (e.g., `Payment`).

### 1. Create Module Structure
Create the following directory structure:
```
internal/payment/
├── domain/
├── adapters/
│   ├── postgres/
│   ├── http/
│   └── grpc/
```

### 2. Define the Domain
In `internal/payment/domain`:
- **Entity (`entity.go`)**: Define core data structures (e.g., `Payment`).
- **Repository Interface (`repository.go`)**: Define `Repository` interface.
- **Service Interface & Implementation (`service.go`)**: Define business logic.

### 3. Implement Adapters
In `internal/payment/adapters`:

#### Database (`internal/payment/adapters/postgres`)
- Implement the `Repository` interface.

#### HTTP (`internal/payment/adapters/http`)
- Create a `Handler` struct.
- Implement methods like `CreatePayment`, `GetPayment`.
- Implement `RegisterRoutes(mux *http.ServeMux)`.

#### gRPC (`internal/payment/adapters/grpc`)
- Define proto in `proto/payment.proto`.
- Generate code.
- Implement the server interface.

### 4. Wiring (`internal/app/wiring.go`)
- Import your new module packages.
- Add Service and Handler to `Application` struct.
- Initialize Repository, Service, and Handler in `Init`.

### 5. Main (`cmd/server/main.go`)
- Call `a.PaymentHTTP.RegisterRoutes(mux)` to expose endpoints.

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
