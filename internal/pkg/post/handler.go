package post

import "net/http"

type Handler interface {
	Details(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
}
