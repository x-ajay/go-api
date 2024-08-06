package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

var (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/postgres?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
