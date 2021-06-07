package thread

import "net/http"

type Handler interface {
	Vote(w http.ResponseWriter, r *http.Request)
	Details(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}
