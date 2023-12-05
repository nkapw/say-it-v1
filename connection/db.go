package connection

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/lib/pq"
)

const (
	host     = "34.128.106.183"
	port     = 5432
	user     = "postgres"
	password = "opH?R+yntz8)n{#M"
	dbname   = "sayit_db"
)

var db *sql.DB

func GetConnection() *sql.DB {
	connectionString :=
		fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
	return db
}
