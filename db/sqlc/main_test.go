package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ftvdexcz/simplebank/config"
	_ "github.com/lib/pq"
)



var testQueries *Queries
var testDB *sql.DB



func TestMain(m *testing.M) {
	var err error
	config, err := config.LoadConfig("../..")
	if err != nil{
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(config.DbDriver, config.DBSource)
	if err != nil{
		log.Fatal("cannot connect to database", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}