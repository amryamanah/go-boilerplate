package store

import (
	"database/sql"
	"fmt"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	err := config.LoadConfig("../../../configs")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Printf("[VIPER] Config: %+v\n", config.Config)

	testDB, err = sql.Open("postgres", config.Config.GetDBConnStr())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
