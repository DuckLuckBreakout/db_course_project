package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	postgreSqlConn, err := sql.Open(
		"postgres",
		"user=forum_root "+
			"password=root "+
			"dbname=forum "+
			"host=localhost "+
			"port= 5432"+
			"sslmode=disable ",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
