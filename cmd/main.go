package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	user_handler "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/handler"
	user_repository "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/repository"
	user_usecase "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/usecase"
)

func main() {
	postgreSqlConn, err := sql.Open(
		"postgres",
		"user=forum_root "+
			"password=root "+
			"dbname=forum "+
			"host=localhost "+
			"port= 5432 "+
			"sslmode=disable ",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()



	userRepository := user_repository.NewRepository(postgreSqlConn)
	userUseCase := user_usecase.NewUseCase(userRepository)
	userHandler := user_handler.NewHandler(userUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/api/user/{nickname:[a-z]+}/create", userHandler.Create).Methods("POST")


	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
