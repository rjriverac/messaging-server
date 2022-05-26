package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/rjriverac/messaging-server/api"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/util"
)



func main() {
	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:",err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db")
	}

	store := db.NewStore(conn)
	server:= api.NewServer(store)

	err = server.StartServer(config.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}
}
