package users

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		err := errors.NewBadRequestError("invalid json")
		c.JSON(err.Status, err)
		return
	}

	result, getErr := services.GetUser(user)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(result.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		err := errors.NewInternalServerError("login failed")
		c.JSON(err.Status, err)
		return
	}

	jwtName := os.Getenv("JWT_NAME")
	jwtPath := os.Getenv("JWT_PATH")
	jwtDomain := os.Getenv("JWT_DOMAIN")
	jwtMaxAgeStr := os.Getenv("JWT_MAXAGE")
	jwtSecureStr := os.Getenv("JWT_SECURE")
	jwtHttpOnlyStr := os.Getenv("JWT_HTTPONLY")

	jwtMaxAge, err := strconv.Atoi(jwtMaxAgeStr)
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	jwtSecure, err := strconv.ParseBool(jwtSecureStr)
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	jwtHttpOnly, err := strconv.ParseBool(jwtHttpOnlyStr)
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	c.SetCookie(jwtName, token, jwtMaxAge, jwtPath, jwtDomain, jwtSecure, jwtHttpOnly)
	c.JSON(http.StatusOK, result)
}
