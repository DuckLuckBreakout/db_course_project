package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/user"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	UseCase    user.UseCase
}

func NewHandler(userUCase user.UseCase) user.Handler {
	return &Handler{
		UseCase:    userUCase,
	}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var newUser models.User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	newUser.Nickname =  mux.Vars(r)["nickname"]

	result, err := h.UseCase.Create(&newUser)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	http_utils.SetJSONResponse(w, result, http.StatusCreated)
}