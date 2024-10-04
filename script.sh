#!/bin/bash

# Get the base directory
BASE_DIR=$(pwd)

echo "Tidying dependencies for the user service"
cd "$BASE_DIR/user-service" || exit
go mod tidy

echo "Tidying dependencies for the product service"
cd "$BASE_DIR/product-service" || exit
go mod tidy

echo "Tidying dependencies for the order service"
cd "$BASE_DIR/order-service" || exit
go mod tidy

echo "Tidying dependencies for the GraphQL gateway"
cd "$BASE_DIR/graphql-gateway" || exit
go mod tidy

echo "All dependencies tidied, starting services..."

# Starting all services after tidying
echo "Starting the user service microservice"
cd "$BASE_DIR/user-service" || exit
go run main.go &
USER_PID=$!

echo "Starting the product service microservice"
cd "$BASE_DIR/product-service" || exit
go run main.go &
PRODUCT_PID=$!

echo "Starting the order service microservice"
cd "$BASE_DIR/order-service" || exit
go run main.go &
ORDER_PID=$!

echo "Starting the GraphQL gateway"
cd "$BASE_DIR/graphql-gateway" || exit
go run server.go &
GRAPHQL_PID=$!

# Wait for all services to complete
wait $USER_PID $PRODUCT_PID $ORDER_PID $GRAPHQL_PID

echo "All services started"
