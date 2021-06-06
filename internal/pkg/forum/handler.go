package forum

import "net/http"

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Details(w http.ResponseWriter, r *http.Request)
	CreateThread(w http.ResponseWriter, r *http.Request)
	Threads(w http.ResponseWriter, r *http.Request)
}
