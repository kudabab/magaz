package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kudabab/market-s/controller"
	"github.com/kudabab/market-s/db"
	"github.com/kudabab/market-s/middleware"
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()
	db.InitDB()

	router := gin.Default()
	router.POST("/register", controller.RegisterUser)

	/*protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())*/
	/*router.GET("/auth", middleware.AuthMiddleware(), func(c *gin.Context) {
		username, err := c.Get("username")
		fmt.Println("Username: ", username)
		if !err {
			c.JSON(500, gin.H{"error": "Username not found"})
			return
		}
		c.JSON(200, gin.H{"message": "good"})
	})*/

	/*router.GET("/auth", middleware.AuthMiddleware(), controller.GetUserByUsername)*/
	router.Use(middleware.AuthMiddleware(), controller.GetUserByUsername)

	router.GET("/users", controller.FindAllUsers)

	router.GET("/items", controller.GetItemsByCategory)  //http://localhost:8080/items?category=water(пример запроса)
	router.GET("/items/item", controller.GetProductByID) //http://localhost:8080/items/item?category=water&id=2(пример запроса)

	router.GET("cart/add", controller.AddToCart) //http://localhost:8080/cart/add?product_id=2&category=water&quantity=1(пример запроса)
	router.POST("cart/remove", controller.RemoveFromCart)
	router.GET("cart", controller.GetCart)
	router.POST("order", controller.CreateOrder)
	router.Run(":8080")
}
