package users

import (
	"net/http"

	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user users.User

	// passing c.Request.Body into the memory address pointed by the pointer
	if err := c.ShouldBindJSON(&user); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		c.JSON(err.Status, err)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusOK, result)
}
