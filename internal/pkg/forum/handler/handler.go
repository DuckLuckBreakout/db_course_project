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
	"strconv"
	"time"
)

type Handler struct {
	UseCase forum.UseCase
}

func (h Handler) Users(w http.ResponseWriter, r *http.Request) {
	var userSearch models.UserSearch

	userSearch.Forum = mux.Vars(r)["slug"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	userSearch.Limit = int32(limit)

	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	userSearch.Desc = desc

	userSearch.Since = r.URL.Query().Get("since")

	result, err := h.UseCase.Users(&userSearch)

	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, result, http.StatusOK)
}

func (h Handler) Threads(w http.ResponseWriter, r *http.Request) {

	var newThreadSearch models.ThreadSearch

	newThreadSearch.Forum = mux.Vars(r)["slug"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	newThreadSearch.Limit = int32(limit)

	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	newThreadSearch.Desc = desc

	sinceString := r.URL.Query().Get("since")
	if sinceString != "" {
		since, _ := time.Parse("2006-01-02T15:04:05.000Z", sinceString)
		newThreadSearch.Since = since
	} else {
		since, _ := time.Parse("2006-01-02T15:04:05.000Z", "3006-01-02T15:04:05.000Z")
		newThreadSearch.Since = since
	}

	result, err := h.UseCase.Threads(&newThreadSearch, sinceString)
	if err == errors.ErrThreadAlreadyCreatedError {
		http_utils.SetJSONResponse(w, newThreadSearch, http.StatusConflict)
		return
	}
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, result, http.StatusOK)
}

func (h Handler) CreateThread(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var newThread models.Thread
	err = json.Unmarshal(body, &newThread)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	newThread.Forum = mux.Vars(r)["slug"]

	err = h.UseCase.CreateThread(&newThread)
	if err == errors.ErrThreadAlreadyCreatedError {
		http_utils.SetJSONResponse(w, newThread, http.StatusConflict)
		return
	}
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, newThread, http.StatusCreated)
}

func (h Handler) Details(w http.ResponseWriter, r *http.Request) {

	var forumForDetails models.Forum

	forumForDetails.Slug = mux.Vars(r)["slug"]

	err := h.UseCase.Details(&forumForDetails)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	http_utils.SetJSONResponse(w, forumForDetails, http.StatusOK)
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
