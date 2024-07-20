package error

import "fmt"

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *RestErr) Error() string {
	return fmt.Sprintf("status %d: %s", e.Status, e.Message)
}
