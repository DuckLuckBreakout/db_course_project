package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	UseCase thread.UseCase
}

func (h Handler) Posts(w http.ResponseWriter, r *http.Request) {

	var newPostSearch models.PostSearch

	slugOrId := mux.Vars(r)["slug_or_id"]
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		newPostSearch.ThreadSlug = slugOrId
	} else {
		newPostSearch.Thread = int32(id)
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	newPostSearch.Limit = int32(limit)

	desc, _ := strconv.ParseBool(r.URL.Query().Get("desc"))
	newPostSearch.Desc = desc

	sort := r.URL.Query().Get("sort")
	newPostSearch.Sort = sort

	since, _ := strconv.Atoi(r.URL.Query().Get("since"))
	newPostSearch.Since = int64(since)
	result, err := h.UseCase.Posts(&newPostSearch)
	if err == errors.ErrThreadAlreadyCreatedError {
		http_utils.SetJSONResponse(w, newPostSearch, http.StatusConflict)
		return
	}
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, result, http.StatusOK)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	posts := make([]*models.Post, 0)
	err = json.Unmarshal(body, &posts)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	slugOrId := mux.Vars(r)["slug_or_id"]

	err = h.UseCase.Create(slugOrId, posts)
	if err == errors.ErrUserNotFound {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	if err == errors.ErrUserAlreadyCreatedError {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}
	if err != nil {
		if err.(*pq.Error).Code == "P0001" {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
			return
		}
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, posts, http.StatusCreated)
}

func (h Handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var threadInfo models.ThreadUpdate
	err = json.Unmarshal(body, &threadInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	slugOrId := mux.Vars(r)["slug_or_id"]
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		threadInfo.Slug = slugOrId
	} else {
		threadInfo.Id = int32(id)
	}

	result, err := h.UseCase.UpdateDetails(&threadInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	http_utils.SetJSONResponse(w, result, http.StatusOK)
}

func (h Handler) Details(w http.ResponseWriter, r *http.Request) {

	var threadInfo models.Thread

	slugOrId := mux.Vars(r)["slug_or_id"]
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		threadInfo.Slug = slugOrId
	} else {
		threadInfo.Id = int32(id)
	}

	err = h.UseCase.Details(&threadInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, threadInfo, http.StatusOK)
}

func (h Handler) Vote(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var voteInfo models.ThreadVoice
	err = json.Unmarshal(body, &voteInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	slugOrId := mux.Vars(r)["slug_or_id"]
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		voteInfo.ThreadSlug = slugOrId
	} else {
		voteInfo.ThreadID = int32(id)
	}

	result, err := h.UseCase.Vote(&voteInfo)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}

	http_utils.SetJSONResponse(w, result, http.StatusOK)
}

func NewHandler(useCase thread.UseCase) thread.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
