package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"whosdriving-be/graph"
	"whosdriving-be/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultHost = "127.0.0.1"
const defaultPort = "8080"
const defaultDbHostPath = "/app/data/whosdriving"

type Config struct {
	host   string
	port   string
	dbPath string
}

func createConfig() Config {
	var config Config
	var found bool
	config.host, found = os.LookupEnv("HOST")
	if !found {
		config.host = defaultHost
	}

	config.port, found = os.LookupEnv("PORT")
	if !found {
		config.port = defaultPort
	}

	config.dbPath, found = os.LookupEnv("DB_PATH")
	if !found {
		config.dbPath = defaultDbHostPath
	}
	return config
}

func newDb(dbPath string) *sql.DB {
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		db.Close()
		log.Fatal(dbErr)
	}

	// and validate dsn
	pingErr := db.Ping()
	if pingErr != nil {
		db.Close()
		log.Fatal(pingErr)
	}

	return db
}

func main() {

	config := createConfig()
	db := newDb(config.dbPath)
	defer db.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://%s:%s/ for GraphQL playground", config.host, config.port)
	log.Fatal(http.ListenAndServe(config.host+":"+config.port, nil))
}
