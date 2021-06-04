package service

import "net/http"

type Handler interface {
	Clear(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
}
