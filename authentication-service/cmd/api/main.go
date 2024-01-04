package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service on port", webPort)
	conn := connectToDB()
	if conn == nil {
		log.Panic("Unable to connect to database. Exiting...")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to database. Retrying...")
			counts++
		} else {
			log.Println("Successfully connected to database.")
			return connection
		}

		if counts > 10 {
			log.Println("Unable to connect to database after 10 tries. Exiting...")
			return nil
		}

		log.Println("Backing off for 5 seconds...")
		time.Sleep(5 * time.Second)
		continue
	}
}
