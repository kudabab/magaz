package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kudabab/market-s/entity"
	"github.com/kudabab/market-s/service"
)

func AddToCart(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(entity.User)
	userID := currentUser.Id

	fmt.Println("userID: ", userID)

	productID, err := strconv.Atoi(c.Query("product_id"))
	fmt.Println("productID: ", productID)
	category := c.Query("category")
	fmt.Println("category:", category)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	quantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil || quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	cart, err := service.AddToCart(c.Request.Context(), userID, productID, quantity, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func RemoveFromCart(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(entity.User)
	userID := currentUser.Id

	productID, _ := strconv.Atoi(c.Query("product_id"))

	cart, err := service.RemoveFromCart(userID, productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// возвращает содержимое корзины
func GetCart(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(entity.User)
	userID := currentUser.Id

	cart, err := service.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// создает заказ на основе корзины
func CreateOrder(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(entity.User)
	userID := currentUser.Id

	order, err := service.CreateOrder(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
