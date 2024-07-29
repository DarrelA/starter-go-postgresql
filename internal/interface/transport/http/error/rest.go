package error

import (
	"net/http"

	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
)

func NewInternalServerError(message string) *restDomainErr.RestErr {
	return &restDomainErr.RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *restDomainErr.RestErr {
	return &restDomainErr.RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewUnprocessableEntityError(message string) *restDomainErr.RestErr {
	return &restDomainErr.RestErr{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
	}
}

func NewUnauthorizedError(message string) *restDomainErr.RestErr {
	return &restDomainErr.RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewBadGatewayError(message string) *restDomainErr.RestErr {
	return &restDomainErr.RestErr{
		Message: message,
		Status:  http.StatusBadGateway,
	}
}
