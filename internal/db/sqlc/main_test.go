package db

import (
	"database/sql"
	"github.com/amryamanah/go-boilerplate/pkg/application"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("failed to load env vars")
	}

	app, err := application.Get(application.TEST)
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(app.DBConn.Client)
	testDB = app.DBConn.Client

	os.Exit(m.Run())
}
