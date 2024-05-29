package response

import (
	"context"
	"fmt"
	"github.com/gmlalfjr/go-service-kit/env"
	"github.com/gmlalfjr/go-service-kit/errs"
	"github.com/gmlalfjr/go-service-kit/transform"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type response struct {
	statusCode       int                `json:"-"`
	Code             int                `json:"code"`
	Message          string             `json:"message"`
	Data             interface{}        `json:"data"`
	ErrorValidations []errorValidations `json:"error_validations,omitempty"`
}

type errorValidations struct {
	Field      string `json:"field"`
	Validation string `json:"validation"`
	Message    string `json:"message"`
}

type Response interface {
	JSON(c *gin.Context)
}

func (r *response) JSON(c *gin.Context) {
	// NOTES: trash way to set http response body
	body, _ := transform.InterfaceToString(r)
	// check response
	if len(body) > env.GetInt("MAX_BODY_SIZE", 1500) {
		body = "success request"
	}

	c.Set("HTTP_RESPONSE_BODY", body)
	if r.statusCode >= http.StatusBadRequest {
		c.AbortWithStatusJSON(r.statusCode, r)
		return
	}

	c.JSON(r.statusCode, r)
}

// Error returns an error response with the given status code and error message
func Error(ctx context.Context, err error) Response {
	var resp = &response{
		statusCode: http.StatusInternalServerError,
		Code:       errs.SOMETHING_WENT_WRONG.Code(),
		Message:    errs.SomethingWentWrong,
	}

	if err != nil {
		switch er := err.(type) {
		case errs.CodeErr:
			resp = &response{
				statusCode: er.StatusCode(),
				Code:       er.Code(),
				Message:    er.Message(),
			}
		case *errs.Error:
			resp = &response{
				statusCode: er.StatusCode(),
				Code:       er.SystemCode(),
				Message:    er.Message(),
			}
		case validator.ValidationErrors:
			resp = &response{
				statusCode:       errs.VALIDATION_ERROR.StatusCode(),
				Code:             errs.VALIDATION_ERROR.Code(),
				Message:          errs.VALIDATION_ERROR.Message(),
				ErrorValidations: make([]errorValidations, 0),
			}

			for _, fe := range er {
				v := fe.Tag()
				if fe.Param() != "" {
					v += fmt.Sprintf("=%s", fe.Param())
				}

				resp.ErrorValidations = append(resp.ErrorValidations, errorValidations{
					Field:      fe.Namespace(),
					Validation: v,
					Message:    fe.Error(),
				})
			}
		}
	}

	return resp
}

// Success returns an success response
func Success(ctx context.Context, statusCode int, data interface{}) Response {
	var successCode CodeSuccess
	switch statusCode {
	case http.StatusOK:
		successCode = SUCCESS_GET
	case http.StatusCreated:
		successCode = SUCCESS_CREATED
	default:
		successCode = SUCCESS_GET
	}

	return &response{statusCode: statusCode, Code: successCode.Code(), Message: successCode.Message(), Data: data}
}
