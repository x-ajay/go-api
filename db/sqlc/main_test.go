package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	utilsconfig "github.com/x-ajay/go-api/utils/config"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utilsconfig.LoadConfig("../../")
	if err != nil {
		panic(err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
