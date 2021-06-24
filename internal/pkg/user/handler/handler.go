package handler

import (
	"encoding/json"
	"fmt"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/user"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	UseCase user.UseCase
}

func NewHandler(userUCase user.UseCase) user.Handler {
	return &Handler{
		UseCase: userUCase,
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

	newUser.Nickname = mux.Vars(r)["nickname"]

	result, err := h.UseCase.Create(&newUser)
	if err == errors.ErrUserAlreadyCreatedError {
		fmt.Println(err)
		http_utils.SetJSONResponse(w, result, http.StatusConflict)
		return
	}

	if err != nil {
		fmt.Println(err)
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	http_utils.SetJSONResponse(w, result[0], http.StatusCreated)
}

func (h Handler) Profile(w http.ResponseWriter, r *http.Request) {

	var userInfo models.User

	userInfo.Nickname = mux.Vars(r)["nickname"]

	err := h.UseCase.Profile(&userInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, userInfo, http.StatusOK)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var userInfo models.User
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	userInfo.Nickname = mux.Vars(r)["nickname"]
	if userInfo.Email == "" && userInfo.Fullname == "" && userInfo.About == "" {
		err = h.UseCase.Profile(&userInfo)
	} else {
		err = h.UseCase.UpdateProfile(&userInfo)
		if err == nil {
			err = h.UseCase.Profile(&userInfo)
		}
	}
	if err == errors.ErrUserNotFound {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}

	http_utils.SetJSONResponse(w, userInfo, http.StatusOK)
}
