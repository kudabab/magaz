package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kudabab/market-s/entity"
	"github.com/kudabab/market-s/service"
)

func RegisterUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenString, err := service.CreateUser(c.Request.Context(), user)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": tokenString})
}

func FindAllUsers(c *gin.Context) {
	users, err := service.FindAllUsers(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error users"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"users": users})
}

func GetUserByUsername(c *gin.Context) {
	user, err := c.Get("user")
	if !err {
		c.IndentedJSON(500, gin.H{"error": "Could not retrieve user information"})
		return
	}

	currentUser := user.(entity.User)

	c.JSON(200, gin.H{
		"usernmae:": currentUser.Username,
		"email":     currentUser.Email,
	})
}
