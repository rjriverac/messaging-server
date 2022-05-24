package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5588/message_db?sslmode=disable"
)

var testQueries *Queries

var testDB *sql.DB
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db")
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
