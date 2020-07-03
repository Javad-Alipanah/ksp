package main

import (
	"ksp/graph"
	"ksp/graph/generated"
	database "ksp/internal/pkg/db/mysql"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = defaultPort
	}

	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")

	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPass == "" {
		log.Fatal("database config not found")
	}

	_ = database.InitDB(dbHost, dbPort, dbName, dbUser, dbPass)
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
