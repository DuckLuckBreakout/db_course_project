package errors

import (
	"fmt"
)

type Error struct {
	Message string `json:"message"`
}

func (err Error) Error() string {
	return fmt.Sprintf("error: happened %s", err.Message)
}

func CreateError(err error) error {
	if _, ok := err.(Error); ok {
		return err
	}

	return Error{Message: err.Error()}
}

var (
	ErrUserNotFound error = Error{
		Message: "Can't find user with id #42\n",
	}
	ErrUserAlreadyCreatedError error = Error{
		Message: "something went wrong",
	}
	ErrForumAlreadyCreatedError error = Error{
		Message: "something went wrong",
	}
	ErrThreadAlreadyCreatedError error = Error{
		Message: "something went wrong",
	}
	ErrBadRequest error = Error{
		Message: "incorrect request",
	}
)
