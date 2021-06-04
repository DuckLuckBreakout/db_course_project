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
	ErrSessionNotFound error = Error{
		Message: "session not found",
	}
	ErrEmailAlreadyExist error = Error{
		Message: "user email already exist",
	}
	ErrServerSystem error = Error{
		Message: "system error",
	}
	ErrFileNotRead error = Error{
		Message: "can't read file",
	}
	ErrIncorrectFileType error = Error{
		Message: "incorrect file type",
	}
	ErrProductNotFound error = Error{
		Message: "product not found",
	}
	ErrCategoryNotFound error = Error{
		Message: "category not found",
	}
	ErrIncorrectPaginator error = Error{
		Message: "incorrect params of pagination",
	}
	ErrBadRequest error = Error{
		Message: "incorrect request",
	}
	ErrCanNotUnmarshal error = Error{
		Message: "can't unmarshal",
	}
	ErrCanNotMarshal error = Error{
		Message: "can't marshal",
	}
	ErrDBInternalError error = Error{
		Message: "internal db error",
	}
	ErrDBFailedConnection error = Error{
		Message: "can't connect to db",
	}
	ErrCartNotFound error = Error{
		Message: "user cart not found",
	}
	ErrProductNotFoundInCart error = Error{
		Message: "product not found in cart",
	}
	ErrInvalidData error = Error{
		Message: "invalid data",
	}
	ErrRequireIdNotFound error = Error{
		Message: "require id not found",
	}
	ErrOpenFile error = Error{
		Message: "can't open file",
	}
	ErrNotFoundCsrfToken error = Error{
		Message: "csrf token not found",
	}
	ErrIncorrectJwtToken error = Error{
		Message: "incorrect jwt token",
	}
	ErrS3InternalError error = Error{
		Message: "can't upload file to S3",
	}
	ErrIncorrectAuthData error = Error{
		Message: "incorrect auth user data",
	}
	ErrNoWriteRights error = Error{
		Message: "no write rights",
	}
	ErrCanNotAddReview error = Error{
		Message: "user can not add review",
	}
	ErrIncorrectSearchQuery error = Error{
		Message: "incorrect search query",
	}
	ErrPromoCodeNotFound error = Error{
		Message: "promo code not found",
	}
	ErrProductNotInPromo error = Error{
		Message: "product does not participate in the promotion",
	}
)
