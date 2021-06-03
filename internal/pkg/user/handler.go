package user

import "net/http"

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
}
