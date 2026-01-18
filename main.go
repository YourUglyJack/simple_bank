package main

import (
	"database/sql"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	_ "github.com/lib/pq"
)

// var Queries *Queries
var DB *sql.DB

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config, err:", err)
	}
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
