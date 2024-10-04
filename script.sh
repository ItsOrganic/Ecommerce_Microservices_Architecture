#!/bin/bash

# Function to check if a command exists
check_command() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a service is running
check_service() {
    if ! systemctl is-active --quiet "$1"; then
        echo "Error: $1 is not running. Please start it before running this script."
        exit 1
    fi
}

# Function to install MongoDB using Docker
install_mongodb() {
    echo "Installing MongoDB using Docker..."
    docker run --name mongodb -d -p 27017:27017 mongo:latest
    echo "MongoDB installed and running on port 27017."
}

# Function to install RabbitMQ using Docker
install_rabbitmq() {
    echo "Installing RabbitMQ using Docker..."
    docker run --name rabbitmq -d -p 5672:5672 -p 15672:15672 rabbitmq:management
    echo "RabbitMQ installed and running on ports 5672 (AMQP) and 15672 (Management UI)."
}

# Function to install Redis using Docker
install_redis() {
    echo "Installing Redis using Docker..."
    docker run --name redis -d -p 6379:6379 redis:latest
    echo "Redis installed and running on port 6379."
}

# Check for Go installation
if ! check_command "go"; then
    echo "Go is not installed. Installing Go using Docker..."
    echo "Please install Go manually, or run the following command:"
    echo "docker run --rm -it golang:latest bash"
    exit 1
fi

# Check for MongoDB installation
if ! check_command "mongo"; then
    install_mongodb
else
    # Check if MongoDB is running
    check_service "mongod"
fi

# Check for RabbitMQ installation
if ! check_command "rabbitmq-server"; then
    install_rabbitmq
else
    # Check if RabbitMQ is running
    check_service "rabbitmq-server"
fi

# Check for Redis installation
if ! check_command "redis-cli"; then
    install_redis
else
    # Check if Redis is running
    if ! redis-cli ping | grep -q "PONG"; then
        echo "Error: Redis is not running. Installing Redis using Docker..."
        install_redis
    else
        echo "Redis is already running."
    fi
fi

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
