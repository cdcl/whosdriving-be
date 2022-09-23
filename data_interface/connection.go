package data_interface

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"
)

type LuwContext struct {
	Conn *sql.DB
	Tx   *sql.Tx
}

func NewConnection(dbPath string) (*sql.DB, error) {
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		return nil, dbErr
	}

	// and validate dsn
	pingErr := db.Ping()
	if pingErr != nil {
		db.Close()
		return nil, pingErr
	}

	return db, nil
}

func Migrate(ddlPath string, db *sql.DB) error {
	file, err := ioutil.ReadFile(ddlPath)
	if err != nil {
		log.Printf("Error: Could not read the file - %s", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error: Couldn't start txn - %s", err)
		return err
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
			log.Printf("Error: Could not execute command [%s] - %s", command, err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error: Could not commit - %s", err)
		return err
	}

	return nil
}
