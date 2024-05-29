package errs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	SomethingWentWrong = "Something went wrong"
	BadRequest         = "Bad request error"
	NotFound           = "Not found"
)

type CodeErr int

const (
	SOMETHING_WENT_WRONG CodeErr = 500
	VALIDATION_ERROR     CodeErr = 600
)

var mapCodeErrStatusCode = map[CodeErr]int{
	SOMETHING_WENT_WRONG: http.StatusInternalServerError,
	VALIDATION_ERROR:     http.StatusBadRequest,
}

var mapCodeErrMessage = map[CodeErr]string{
	SOMETHING_WENT_WRONG: "Something went wrong",
	VALIDATION_ERROR:     "Error Validation",
}

func (ce CodeErr) Error() string {
	return fmt.Sprint(strings.ToLower(ce.Message()))
}

func (ce CodeErr) Errors() error {
	return NewError(errors.New(strings.ToLower(ce.Message())), ce.StatusCode(), ce.Code(), ce.Message())
}

func (ce CodeErr) Code() int {
	return int(ce)
}

func (ce CodeErr) StatusCode() int {
	val, ok := mapCodeErrStatusCode[ce]
	if !ok {
		return http.StatusInternalServerError
	}

	return val
}

func (ce CodeErr) Message() string {
	val, ok := mapCodeErrMessage[ce]
	if !ok {
		return SomethingWentWrong
	}

	return val
}
