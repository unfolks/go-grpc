#!/bin/bash

# gRPC Testing Script for Product, Customer, and Order Services
# Make sure your server is running on localhost:50051 before running this script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Server address
SERVER_ADDRESS="localhost:50051"

echo -e "${YELLOW}gRPC Testing Script${NC}"
echo -e "${YELLOW}====================${NC}"
echo -e "Testing server at: ${SERVER_ADDRESS}"
echo ""

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}Error: grpcurl is not installed. Please install it first:${NC}"
    echo "brew install grpcurl"
    exit 1
fi

# Function to test a gRPC call
test_grpc_call() {
    local service=$1
    local method=$2
    local data=$3
    local description=$4
    
    echo -e "${GREEN}Testing: $description${NC}"
    echo "Command: grpcurl -plaintext -d '$data' $SERVER_ADDRESS $service/$method"
    
    grpcurl -plaintext -d "$data" "$SERVER_ADDRESS" "$service/$method"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
    fi
    echo "----------------------------------------"
    echo ""
}

# List all available services
echo -e "${YELLOW}Available Services:${NC}"
grpcurl -plaintext "$SERVER_ADDRESS" list
echo "----------------------------------------"
echo ""

# ===== PRODUCT SERVICE TESTS =====
echo -e "${YELLOW}===== PRODUCT SERVICE TESTS =====${NC}"

# Create Product
test_grpc_call "productpb.ProductService" "CreateProduct" \
'{"name": "Laptop", "price": 999.99}' \
"Create Product"

# Create another product
test_grpc_call "productpb.ProductService" "CreateProduct" \
'{"name": "Mouse", "price": 29.99}' \
"Create Another Product"

# List Products
test_grpc_call "productpb.ProductService" "ListProducts" \
'{}' \
"List All Products"

# Get Product (assuming ID 1 exists)
test_grpc_call "productpb.ProductService" "GetProduct" \
'{"id": "1"}' \
"Get Product with ID 1"

# Update Product (assuming ID 1 exists)
test_grpc_call "productpb.ProductService" "UpdateProduct" \
'{"id": "1", "name": "Updated Laptop", "price": 1099.99}' \
"Update Product with ID 1"

# Delete Product (assuming ID 2 exists)
test_grpc_call "productpb.ProductService" "DeleteProduct" \
'{"id": "2"}' \
"Delete Product with ID 2"

# ===== ORDER SERVICE TESTS =====
echo -e "${YELLOW}===== ORDER SERVICE TESTS =====${NC}"

# Create Order
test_grpc_call "orderpb.ORderService" "CreateOrder" \
'{"amount": 1299.98}' \
"Create Order"

# Create another order
test_grpc_call "orderpb.ORderService" "CreateOrder" \
'{"amount": 59.99}' \
"Create Another Order"

# List Orders
test_grpc_call "orderpb.ORderService" "ListOrders" \
'{}' \
"List All Orders"

# Get Order (assuming ID 1 exists)
test_grpc_call "orderpb.ORderService" "GetOrder" \
'{"id": "1"}' \
"Get Order with ID 1"

# Update Order (assuming ID 1 exists)
test_grpc_call "orderpb.ORderService" "UpdateOrder" \
'{"id": "1", "amount": 1399.99}' \
"Update Order with ID 1"

# Delete Order (assuming ID 2 exists)
test_grpc_call "orderpb.ORderService" "DeleteOrder" \
'{"id": "2"}' \
"Delete Order with ID 2"

# ===== CUSTOMER SERVICE TESTS =====
echo -e "${YELLOW}===== CUSTOMER SERVICE TESTS =====${NC}"

# Create Customer
test_grpc_call "customerpb.CustomerService" "CreateCustomer" \
'{"id": "1", "name": "John Doe", "email": "john.doe@example.com", "address": "123 Main St, New York, NY 10001"}' \
"Create Customer"

# Create another customer
test_grpc_call "customerpb.CustomerService" "CreateCustomer" \
'{"id": "2", "name": "Jane Smith", "email": "jane.smith@example.com", "address": "456 Oak Ave, Los Angeles, CA 90001"}' \
"Create Another Customer"

echo -e "${GREEN}gRPC Testing Complete!${NC}"
echo -e "${YELLOW}Note: Based on the proto files, the customer service only has a CreateCustomer method.${NC}"
echo -e "${YELLOW}Note: Some tests assume certain IDs exist. Adjust IDs as needed based on your data.${NC}"