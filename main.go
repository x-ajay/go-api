package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/x-ajay/go-api/api"
	db "github.com/x-ajay/go-api/db/sqlc"
	utilsconfig "github.com/x-ajay/go-api/utils/config"
)

func main() {
	config, err := utilsconfig.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBDriver)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.SetupServer(store)

	if err = server.Start(config.HTTPAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

}
