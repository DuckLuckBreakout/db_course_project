package main

import (
	"database/sql"
	forum_handler "github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum/handler"
	forum_repository "github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum/repository"
	forum_usecase "github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum/usecase"
	service_handler "github.com/DuckLuckBreakout/db_course_project/internal/pkg/service/handler"
	service_repository "github.com/DuckLuckBreakout/db_course_project/internal/pkg/service/repository"
	service_usecase "github.com/DuckLuckBreakout/db_course_project/internal/pkg/service/usecase"
	thread_handler "github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread/handler"
	thread_repository "github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread/repository"
	thread_usecase "github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread/usecase"
	user_handler "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/handler"
	user_repository "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/repository"
	user_usecase "github.com/DuckLuckBreakout/db_course_project/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
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
			"port= 5434 "+
			"sslmode=disable ",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()

	userRepository := user_repository.NewRepository(postgreSqlConn)
	userUseCase := user_usecase.NewUseCase(userRepository)
	userHandler := user_handler.NewHandler(userUseCase)

	serviceRepository := service_repository.NewRepository(postgreSqlConn)
	serviceUseCase := service_usecase.NewUseCase(serviceRepository)
	serviceHandler := service_handler.NewHandler(serviceUseCase)

	forumRepository := forum_repository.NewRepository(postgreSqlConn)
	forumUseCase := forum_usecase.NewUseCase(forumRepository)
	forumHandler := forum_handler.NewHandler(forumUseCase)

	threadRepository := thread_repository.NewRepository(postgreSqlConn)
	threadUseCase := thread_usecase.NewUseCase(threadRepository)
	threadHandler := thread_handler.NewHandler(threadUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/api/user/{nickname}/create", userHandler.Create).Methods("POST")
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.Profile).Methods("GET")
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.Update).Methods("POST")

	router.HandleFunc("/api/service/clear", serviceHandler.Clear).Methods("POST")
	router.HandleFunc("/api/service/status", serviceHandler.Status).Methods("GET")

	router.HandleFunc("/api/forum/create", forumHandler.Create).Methods("POST")
	router.HandleFunc("/api/forum/{slug}/details", forumHandler.Details).Methods("GET")
	router.HandleFunc("/api/forum/{slug}/create", forumHandler.CreateThread).Methods("POST")
	router.HandleFunc("/api/forum/{slug}/threads", forumHandler.Threads).Methods("GET")

	router.HandleFunc("/api/thread/{slug_or_id}/vote", threadHandler.Vote).Methods("POST")
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.Details).Methods("GET")
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.UpdateDetails).Methods("POST")
	router.HandleFunc("/api/thread/{slug_or_id}/create", threadHandler.Create).Methods("POST")

	server := &http.Server{
		Addr:         ":5000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
