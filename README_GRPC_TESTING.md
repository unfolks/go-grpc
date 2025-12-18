# gRPC Testing Guide

This guide explains how to test the gRPC services for your Go application.

## Prerequisites

1. **Install grpcurl** - A command-line tool for interacting with gRPC servers
   ```bash
   # On macOS
   brew install grpcurl
   
   # On other platforms, follow: https://github.com/fullstorydev/grpcurl#installation
   ```

2. **Start your server** - Make sure your gRPC server is running on localhost:50051
   ```bash
   go run cmd/server/main.go
   ```

## Testing Methods

### Method 1: Using the Automated Script (Recommended)

Run the automated test script:
```bash
./grpc_test.sh
```

This script will test all available gRPC endpoints with sample data and provide colored output for success/failure.

### Method 2: Manual Testing with grpcurl

You can also test individual endpoints manually:

#### List Available Services
```bash
grpcurl -plaintext localhost:50051 list
```

#### Product Service Tests

**Create a Product:**
```bash
grpcurl -plaintext -d '{"name": "Laptop", "price": 999.99}' \
  localhost:50051 productpb.ProductService/CreateProduct
```

**List All Products:**
```bash
grpcurl -plaintext -d '{}' \
  localhost:50051 productpb.ProductService/ListProducts
```

**Get a Specific Product:**
```bash
grpcurl -plaintext -d '{"id": "1"}' \
  localhost:50051 productpb.ProductService/GetProduct
```

**Update a Product:**
```bash
grpcurl -plaintext -d '{"id": "1", "name": "Updated Laptop", "price": 1099.99}' \
  localhost:50051 productpb.ProductService/UpdateProduct
```

**Delete a Product:**
```bash
grpcurl -plaintext -d '{"id": "1"}' \
  localhost:50051 productpb.ProductService/DeleteProduct
```

#### Order Service Tests

**Create an Order:**
```bash
grpcurl -plaintext -d '{"amount": 1299.98}' \
  localhost:50051 orderpb.ORderService/CreateOrder
```

**List All Orders:**
```bash
grpcurl -plaintext -d '{}' \
  localhost:50051 orderpb.ORderService/ListOrders
```

**Get a Specific Order:**
```bash
grpcurl -plaintext -d '{"id": "1"}' \
  localhost:50051 orderpb.ORderService/GetOrder
```

**Update an Order:**
```bash
grpcurl -plaintext -d '{"id": "1", "amount": 1399.99}' \
  localhost:50051 orderpb.ORderService/UpdateOrder
```

**Delete an Order:**
```bash
grpcurl -plaintext -d '{"id": "1"}' \
  localhost:50051 orderpb.ORderService/DeleteOrder
```

#### Customer Service Tests

**Create a Customer:**
```bash
grpcurl -plaintext -d '{"id": "1", "name": "John Doe", "email": "john.doe@example.com", "address": "123 Main St, New York, NY 10001"}' \
  localhost:50051 customerpb.CustomerService/CreateCustomer
```

## Service Limitations

Based on the proto file analysis:
- **Product Service**: Full CRUD operations available
- **Order Service**: Full CRUD operations available  
- **Customer Service**: Only CreateCustomer method is implemented

## Notes

1. The gRPC server runs on port 50051 (as configured in cmd/server/main.go)
2. Some tests assume certain IDs exist. You may need to adjust IDs based on your actual data
3. The automated script provides colored output for better readability
4. Make sure your PostgreSQL database is running and accessible with the credentials configured in main.go

## Troubleshooting

If you encounter connection issues:
1. Ensure the server is running: `go run cmd/server/main.go`
2. Check that port 50051 is not blocked by firewall
3. Verify your database connection settings in cmd/server/main.go