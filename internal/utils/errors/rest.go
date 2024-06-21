package errors

import (
	"fmt"
	"net/http"
)

type RestErr struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	ErrorType string `json:"error"`
}

func (e *RestErr) Error() string {
	return fmt.Sprintf("status %d: %s", e.Status, e.Message)
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:   message,
		Status:    http.StatusInternalServerError,
		ErrorType: "internal_server_error",
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:   message,
		Status:    http.StatusBadRequest,
		ErrorType: "bad_request",
	}
}

func NewUnprocessableEntityError(message string) *RestErr {
	return &RestErr{
		Message:   message,
		Status:    http.StatusUnprocessableEntity,
		ErrorType: "unprocessable_entity",
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message:   message,
		Status:    http.StatusForbidden,
		ErrorType: "forbidden",
	}
}
