// controllers/itemController.go
package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kudabab/market-s/service"
)

// GetItemsByCategory - функция для получения товаров по категории
func GetItemsByCategory(c *gin.Context) {

	user, exists := c.Get("user")
	fmt.Printf("user = %v exists = %v", user, exists)
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	category := c.Query("category") // Получаем категорию из запроса

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	// Вызываем функцию сервиса для получения товаров по категории
	items, err := service.GetProduct(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get items"})
		return
	}

	// Возвращаем список товаров клиенту
	c.JSON(http.StatusOK, items)
}

func GetProductByID(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	category := c.Query("category")
	idStr := c.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	product, err := service.GetProductByID(category, id)
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}
