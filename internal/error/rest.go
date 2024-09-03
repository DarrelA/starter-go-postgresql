package error

import (
	"fmt"
	"net/http"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *RestErr) Error() string {
	return fmt.Sprintf("status %d: %s", e.Status, e.Message)
}

// Factory functions
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewUnprocessableEntityError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewBadGatewayError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadGateway,
	}
}
