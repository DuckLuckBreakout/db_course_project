package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	UseCase forum.UseCase
}

func (h Handler) Details(w http.ResponseWriter, r *http.Request) {
	var forum models.Forum

	forum.Slug = mux.Vars(r)["slug"]

	err := h.UseCase.Details(&forum)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, forum, http.StatusOK)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var newForum models.Forum
	err = json.Unmarshal(body, &newForum)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	err = h.UseCase.Create(&newForum)
	if err == errors.ErrUserNotFound {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	if err != nil {
		var badResult models.ForumEmpty
		badResult.Slug = newForum.Slug
		badResult.User = newForum.User
		badResult.Title = newForum.Title
		http_utils.SetJSONResponse(w, badResult, http.StatusConflict)
		return
	}

	http_utils.SetJSONResponse(w, newForum, http.StatusCreated)
}

func NewHandler(useCase forum.UseCase) forum.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
