package handler

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
)

type Handler struct {
	UseCase thread.UseCase
}

func NewHandler(useCase thread.UseCase) thread.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
