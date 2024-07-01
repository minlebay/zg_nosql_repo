package model

import "errors"

const (
	NotFound        = "NotFound"
	notFoundMessage = "record not found"

	ResourceAlreadyExists     = "ResourceAlreadyExists"
	alreadyExistsErrorMessage = "resource already exists"

	RepositoryError        = "RepositoryError"
	repositoryErrorMessage = "error in repository operation"

	UnknownError        = "UnknownError"
	unknownErrorMessage = "something went wrong"
)

type AppError struct {
	Err  error
	Type string
}

func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType string) *AppError {
	var err error

	switch errType {
	case NotFound:
		err = errors.New(notFoundMessage)
	case ResourceAlreadyExists:
		err = errors.New(alreadyExistsErrorMessage)
	case RepositoryError:
		err = errors.New(repositoryErrorMessage)
	default:
		err = errors.New(unknownErrorMessage)
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}
