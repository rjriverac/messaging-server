package main

import (
	"database/sql"
	"log"

	"github.com/rjriverac/messaging-server/api"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5588/message_db?sslmode=disable"
	address = "0.0.0.0:8080"
)


func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db")
	}

	store := db.NewStore(conn)
	server:= api.NewServer(store)

	err = server.StartServer(address)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}
}
