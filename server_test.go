package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDb(t *testing.T) {
	os.Remove("./test_new_db.sqlite3")

	db := newDb("./test_new_db.sqlite3")
	defer db.Close()

	_, err := db.Exec("create table foo (id integer not null primary key, name text);")
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO foo(id, name) values(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	var expected = [...]string{"Fatima", "Robert", "Lisa", "Ahmed", "Inconnu"}

	for i, n := range expected {
		_, err = stmt.Exec(i, n)
		if err != nil {
			t.Fatalf("on row %d : %s", i, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo order by id desc")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%d|%s", id, name)
		assert.EqualValues(t, expected[id], name)
	}

	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCustomConfig(t *testing.T) {
	var expectedConfig Config
	expectedConfig.host = "BBBB"
	expectedConfig.port = "AAAA"
	expectedConfig.dbPath = "CCCC"

	os.Setenv("HOST", expectedConfig.host)
	os.Setenv("PORT", expectedConfig.port)
	os.Setenv("DB_PATH", expectedConfig.dbPath)

	config := createConfig()
	t.Log(config)
	assert.EqualValues(t, config.host, expectedConfig.host)
	assert.EqualValues(t, config.port, expectedConfig.port)
	assert.EqualValues(t, config.dbPath, expectedConfig.dbPath)
}

func TestDefaultConfig(t *testing.T) {
	os.Unsetenv("HOST")
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")

	config := createConfig()
	t.Log(config)
	assert.EqualValues(t, config.host, defaultHost)
	assert.EqualValues(t, config.port, defaultPort)
	assert.EqualValues(t, config.dbPath, defaultDbHostPath)
}
