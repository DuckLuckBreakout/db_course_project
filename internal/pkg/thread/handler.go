package thread

import "net/http"

type Handler interface {
	Vote(w http.ResponseWriter, r *http.Request)
}
