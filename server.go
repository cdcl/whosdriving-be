package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"whosdriving-be/data_interface"
	"whosdriving-be/graph"
	"whosdriving-be/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultHost = "127.0.0.1"
const defaultPort = "8080"
const defaultDbHostPath = "/app/data/whosdriving"
const defaultDdlPath = "/app/assets/ddl.whosdriving-core"

type Config struct {
	host    string
	port    string
	dbPath  string
	ddlPath string
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func createConfig() Config {
	log.Printf("Read config from env")
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

	config.ddlPath, found = os.LookupEnv("DDL_PATH")
	if !found {
		config.ddlPath = defaultDdlPath
	}

	log.Printf("%q", config)
	return config
}

func newDb(dbPath string, ddlPath string) *sql.DB {
	newDb := !checkFileExists(dbPath)

	log.Printf("Open database %s", dbPath)
	db, err := data_interface.NewConnection(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Migration this is a new DB
	if newDb && checkFileExists(ddlPath) {
		log.Printf("Migrate database %s", ddlPath)
		err := data_interface.Migrate(ddlPath, db)
		if err != nil {
			db.Close()
			log.Fatal(err)
		}
	}

	return db
}

func main() {
	config := createConfig()
	db := newDb(config.dbPath, config.ddlPath)
	defer db.Close()

	log.Println("Prepare graphQL resolver")
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	log.Println("Setup router")
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Connect to http://%s:%s/ for GraphQL playground", config.host, config.port)
	log.Fatal(http.ListenAndServe(":"+config.port, nil))
}
