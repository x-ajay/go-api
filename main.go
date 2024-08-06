package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/x-ajay/go-api/api"
	db "github.com/x-ajay/go-api/db/sqlc"
)

var (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.SetupServer(store)

	if err = server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

}
