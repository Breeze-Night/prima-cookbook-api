package main

import (
	"net/http"

	"prima_cookbook/config"
	"prima_cookbook/middleware"
	"prima_cookbook/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.InitDB()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", GetHome)

		user := v1.Group("/user")
		{
			user.POST("/register", routes.RegisterUser)
			user.POST("/login", routes.GenerateToken)
		}

		recipe := v1.Group("/recipe")
		{
			recipe.GET("/", routes.GetRecipes)
			recipe.GET("/:id", routes.GetRecipeByID)
			recipe.Use(middleware.Auth()).POST("/", routes.CreateRecipe)
			recipe.Use(middleware.Auth()).PUT("/:id", routes.EditRecipe)
			recipe.Use(middleware.Auth()).DELETE("/:id", routes.DeleteRecipe)
		}

		ingredient := v1.Group("/ingredient")
		{
			ingredient.GET("/", routes.GetIngredients)
			ingredient.GET("/:id", routes.GetIngredientByID)
			ingredient.Use(middleware.Auth()).POST("/:recipe_id", routes.AddIngredientToRecipe)
			ingredient.Use(middleware.Auth()).DELETE("/:recipe_id/:ingredient_id", routes.DeleteIngredientFromRecipe)

		}

		filter := v1.Group("/filter")
		{
			filter.GET("/recipe/ingredient/:ingredient_id", routes.GetRecipesByIngredient)
			filter.GET("/recipe/user/:user_id", routes.GetRecipeByUserID)
		}

	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func GetHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Welcome to Prima Cookbook",
	})
}
