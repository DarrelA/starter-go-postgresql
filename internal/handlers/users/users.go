package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtCfg = configs.JWTSettings

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

	token, err := claims.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		err := errors.NewInternalServerError("login failed")
		c.JSON(err.Status, err)
		return
	}

	c.SetCookie(jwtCfg.Name, token, jwtCfg.MaxAge, jwtCfg.Path, jwtCfg.Domain, jwtCfg.Secure, jwtCfg.HttpOnly)
	c.JSON(http.StatusOK, result)
}

func Get(c *gin.Context) {
	cookie, err := c.Cookie(jwtCfg.Name)
	if err != nil {
		getErr := errors.NewInternalServerError("could not retrieve cookie")
		c.JSON(getErr.Status, getErr)
		return
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(jwtCfg.Secret), nil
	})

	if err != nil {
		restErr := errors.NewInternalServerError("error parsing cookie")
		c.JSON(restErr.Status, restErr)
		return
	}

	// set `claims` variable to be a type of `RegisteredClaims`
	// so we can access `claims.Issuer`
	claims := token.Claims.(*jwt.RegisteredClaims)
	issuer, err := strconv.ParseInt(claims.Issuer, 10, 64)
	if err != nil {
		restErr := errors.NewBadRequestError("user id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	// basically issuer is our user id
	result, restErr := services.GetUserByID(issuer)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func Logout(c *gin.Context) {
	c.SetCookie(jwtCfg.Name, "", -1, "", "", jwtCfg.Secure, jwtCfg.HttpOnly)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
