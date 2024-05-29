package response

import "net/http"

const (
	successMessage = "Sukses"
)

type CodeSuccess int

const (
	SUCCESS_GET     CodeSuccess = 200
	SUCCESS_CREATED CodeSuccess = 201
)

var mapCodeSuccessStatusCode = map[CodeSuccess]int{
	SUCCESS_GET:     http.StatusOK,
	SUCCESS_CREATED: http.StatusCreated,
}

var mapCodeSuccessMessage = map[CodeSuccess]string{
	SUCCESS_GET:     successMessage,
	SUCCESS_CREATED: "Created",
}

func (cs CodeSuccess) StatusCode() int {
	val, ok := mapCodeSuccessStatusCode[cs]
	if !ok {
		return http.StatusOK
	}

	return val
}

func (cs CodeSuccess) Message() string {
	val, ok := mapCodeSuccessMessage[cs]
	if !ok {
		return successMessage
	}

	return val
}

func (cs CodeSuccess) Code() int {
	return int(cs)
}
