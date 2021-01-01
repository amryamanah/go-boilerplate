package domain

import (
	"database/sql"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("failed to load env vars")
	}
	cfg := config.Get()

	var err error
	testDB, err = sql.Open("postgres", cfg.GetTestDBConnStr())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
