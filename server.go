package main

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

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

	// Migration this is a new DB
	if newDb && checkFileExists(ddlPath) {
		log.Printf("Create database with %s", ddlPath)
		file, err := ioutil.ReadFile(ddlPath)
		if err != nil {
			db.Close()
			log.Fatalf("Could not read the file due to this %s error", err)
		}

		tx, err := db.Begin()
		if err != nil {
			db.Close()
			log.Fatalf("Couldn't start txn error: %s", err)
		}
		defer tx.Rollback()

		// Here we searching for ; to split the file in commands
		for _, chunk := range strings.Split(string(file), ";") {
			var commandChunks []string
			// Possibility to add full line comments with # (that we ignore here)
			for _, line := range strings.Split(chunk, "\n") {
				cleanLine := strings.TrimLeft(line, " ")
				// skip empty lines
				if len(cleanLine) == 0 {
					continue
				}

				// Line is not a comment
				if len(cleanLine) < 2 || cleanLine[0:2] != "--" {
					commandChunks = append(commandChunks, cleanLine)
					continue
				}

				// If comment start with #Print display it in console
				if len(cleanLine) > 8 && cleanLine[0:8] == "--Print:" {
					log.Println(strings.TrimSpace(cleanLine[8:]))
				}
			}

			command := strings.Join(commandChunks, "\n")
			//log.Println(command)
			if _, err := tx.Exec(command); err != nil {
				db.Close()
				log.Fatalf("Fatal: Could not execute command [%s] error : %s", command, err)
			}
		}

		if err := tx.Commit(); err != nil {
			db.Close()
			log.Fatalf("Fatal: Could not commit %s", err)
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
	log.Fatal(http.ListenAndServe(config.host+":"+config.port, nil))
}
