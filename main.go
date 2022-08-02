package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/rxbenefits/go-hw/api"
	"log"
	"os"
)

func main() {
	// connection string
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// open database
	psqlConn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	// open database
	db, err := sql.Open("pgx", psqlConn)
	if err != nil {
		panic(err)
	}

	// close database
	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()

	// check db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected!")
	server := api.NewServer(db)
	server.HTTPServe("8080")
}
