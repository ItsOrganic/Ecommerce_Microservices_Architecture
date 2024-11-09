# Microservice Overview
## Table of Contents
- [Requirements](#requirements)
- [Setting Up](#setting-up)
- [Verifying the Setup](#verifying-the-setup)
- [POSTMAN WORKSPACE](#postman-workspace)
- [Accessing the Services](#accessing-the-services)
- [Steps](#steps)
- [Overview of Microservice 1 - User Service](#user-service)
- [Overview of Microservice 2 - Product Service](#product-service)
- [Overview of Microservice 3 - Order Service](#order-service)
- [Overview of GraphQl-Gateway ](#graphql-gateway)
- [Common Components in each service ](#common-components)
- [Overview of GraphQl-Gateway ](#prometheus-configuration)
- [Overview of Script ](#script)



## Requirements 
- Docker is mandatory.
- Golang, RabbitMQ, Redis, Mongo required but if not present it will automatically get downloaded by running the setup script using docker.

  # Setting Up
- Clone the repository and go in the Ecommerce_Microservices_Architecture path.
  <pre><code>git clone https://github.com/ItsOrganic/Ecommerce_Microservices_Architecture.git && cd Ecommerce_Microservices_Architecture</code></pre> 
- Make the script executable and run it will download all the dependencies as well as things required.

  <pre><code>chmod +x script.sh </code></pre> 
- Run the script
   <pre><code>./script.sh</code></pre> 
   If Mongo, Redis, RabbitMq and other things required for running are not present it will download using docker.

  # Verifying the Setup
- **User Service**: Should be running on [http://localhost:8081](http://localhost:8081).
- **Product Service**: Should be running on [http://localhost:8082](http://localhost:8082).
- **Order Service**: Should be running on [http://localhost:8083](http://localhost:8083).
- **GraphQL Gateway**: Should be running on [http://localhost:8080/query](http://localhost:8080/query).
- **RabbitMQ Service**: Should be running on [http://localhost:15672](http://localhost:15672).

# POSTMAN WORKSPACE 

## [https://www.postman.com/itsorganic/assignment/collection/66ff293b3d6feb3fae738e7d/graphql-queries-and-mutations?action=share&creator=28479580](https://www.postman.com/itsorganic/assignment/collection/66ff293b3d6feb3fae738e7d/graphql-queries-and-mutations?action=share&creator=28479580))
---



# Accessing the Services

## User Service [http://localhost:8081](http://localhost:8081).
- **Register User**: `POST /register`
- **Login User**: `POST /login`
- **Get Users**: `GET /users`
- **Get User by ID**: `GET /user/:id`
- **Get Profile by ID**: `GET /profile/:id`
- **Update Profile by ID**: `PUT /profile/:id`
- **Metrics**: `GET /metrics`

## Product Service  [http://localhost:8082](http://localhost:8082)
- **Create Product**: `POST /product`
- **Get Product by Name**: `GET /product/:name`
- **Get Products**: `GET /products`
- **Update Inventory**: `PUT /product/:name`
- **Delete Product**: `DELETE /product/:name`
- **Metrics**: `GET /metrics`

## Order Service  [http://localhost:8083](http://localhost:8083)
- **Create Order**: `POST /order`
- **Get Orders**: `GET /orders`
- **Get Order by ID**: `GET /order/:id`
- **Update Order Status**: `PUT /order/:id`
- **Metrics**: `GET /metrics`

## GraphQL Gateway  [http://localhost:8080](http://localhost:8080)
- **GraphQL Endpoint**: `POST /graphql`
- ** GraphQL Query **: `POST /query`

# User Service

**File**: `user-service/main.go`

## Description
The User Service handles user registration and retrieval. It manages user information such as name and email.

## Endpoints
- **GET /metrics**: Exposes Prometheus metrics.
- **POST /register**: Registers a new user.
- **GET /users**: Retrieves all users.
- **GET /user/:id**: Retrieves a specific user by ID.

## Key Functions
- **Database Connection**: Connects to MongoDB using `db.Connect`.
- **Message Queue Initialization**: Initializes RabbitMQ using `utils.InitMQ` and `utils.CloseMQ`.
- **Metrics**: Initializes Prometheus metrics using `metrics.Init`.

---

# Product Service

**File**: `product-service/main.go`

## Description
The Product Service manages product information, including creation, retrieval, updating, and deletion of products.

## Endpoints
- **GET /metrics**: Exposes Prometheus metrics.
- **POST /product**: Creates a new product.
- **GET /product/:name**: Retrieves a specific product by name.
- **GET /products**: Retrieves all products.
- **PUT /product/:name**: Updates the inventory of a specific product by name.
- **DELETE /product/:name**: Deletes a specific product by name.

## Key Functions
- **Database Connection**: Connects to MongoDB using `db.Connect`.
- **Message Queue Initialization**: Initializes RabbitMQ using `utils.InitMQ` and `utils.CloseMQ`.
- **Metrics**: Initializes Prometheus metrics using `metrics.Init`.

---

# Order Service

**File**: `order-service/main.go`

## Description
The Order Service handles the creation, retrieval, and updating of orders. It interacts with the Product Service to verify product availability and pricing.

## Endpoints
- **GET /metrics**: Exposes Prometheus metrics.
- **GET /orders**: Retrieves all orders.
- **POST /order**: Creates a new order.
- **GET /order/:id**: Retrieves a specific order by ID.
- **PUT /order/:id**: Updates the status of a specific order by ID.

## Key Functions
- **Database Connection**: Connects to MongoDB using `db.Connect`.
- **Message Queue Initialization**: Initializes RabbitMQ using `utils.InitMQ` and `utils.CloseMQ`.
- **Metrics**: Initializes Prometheus metrics using `metrics.Init`.

---

# GraphQL Gateway

**File**: `graphql-gateway/server.go`

## Description
The GraphQL Gateway provides a unified GraphQL API for interacting with the User, Product, and Order services. It aggregates data from these services and exposes it through a single GraphQL endpoint.

## Endpoints
- **POST /graphql**: Handles GraphQL queries and mutations.

## Key Functions
- **Schema Definition**: Defines the GraphQL schema in `schema.graphqls`.
- **Resolvers**: Implements the resolvers for the GraphQL schema in `resolver.go`.
- **Server Initialization**: Initializes the GraphQL server in `server.go`.

---

# Common Components

## Database Connection

**File**: `db/db.go`  
**Description**: Provides functions to connect to MongoDB.

---

## Message Queue Initialization

**File**: `utils/utils.go`  
**Description**: Provides functions to initialize and close RabbitMQ connections.

---

## Metrics

**File**: `metrics/metrics.go`  
**Description**: Provides functions to initialize Prometheus metrics.

---

# Prometheus Configuration

**File**: `user-service/prometheus.yml`

## Description
Configures Prometheus to scrape metrics from the User, Product, and Order services.

## Scrape Configurations
- **User Service**: `localhost:8081`
- **Product Service**: `localhost:8082`
- **Order Service**: `localhost:8083`

---

# Script

**File**: `script.sh`

## Description
A shell script to tidy dependencies and start all microservices.

## Steps
1. Tidies dependencies for each service.
2. Starts the User Service.
3. Starts the Product Service.
4. Starts the Order Service.
5. Starts the GraphQL Gateway.

---
