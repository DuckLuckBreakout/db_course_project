package handler

import (
	"fmt"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"net/http"
)

type Handler struct {
	UseCase service.UseCase
}

func (h Handler) Clear(w http.ResponseWriter, r *http.Request) {
	err := h.UseCase.Clear()
	fmt.Println(err)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	http_utils.SetJSONResponse(w, nil, http.StatusOK)
}

func (h Handler) Status(w http.ResponseWriter, r *http.Request) {
	status, err := h.UseCase.Status()
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	http_utils.SetJSONResponse(w, status, http.StatusOK)
}

func NewHandler(useCase service.UseCase) service.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
