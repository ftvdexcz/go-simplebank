package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/ftvdexcz/simplebank/api"
	"github.com/ftvdexcz/simplebank/config"
	db "github.com/ftvdexcz/simplebank/db/sqlc"
)



func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil{
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil{
		log.Fatal("cannot start server", err)
	}
}