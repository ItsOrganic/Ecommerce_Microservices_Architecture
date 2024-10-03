package main

import (
	"gpql-gateway/db"
	"gpql-gateway/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	err := db.Connect("mongodb://localhost:27017", "user-service", "users")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	err = db.Connect("mongodb://localhost:27017", "product-db", "products")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
