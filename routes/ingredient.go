package routes

import (
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetIngredients(c *gin.Context) {
	ingredients := []models.Ingredient{}
	config.DB.Find(&ingredients)

	config.DB.Preload(clause.Associations).Find(&ingredients)

	c.JSON(http.StatusOK, gin.H{
		"message": "All registered ingredients",
		"data":    ingredients,
	})
}

func GetIngredientByID(c *gin.Context) {
	id := c.Param("id")

	var ingredient models.Ingredient

	data := config.DB.Preload(clause.Associations).First(&ingredient, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Data Not Found",
			"message": "Recipe does not exist",
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Successful",
		"data":    ingredient,
	})
}

func AddIngredientToRecipe(c *gin.Context) {
	var recipe models.Recipe
	var ingredient models.Ingredient
	recipeID := c.Param("recipe_id")

	config.DB.First(&recipe, recipeID)
	err := c.BindJSON(&ingredient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})
		c.Abort()
		return
	}

	config.DB.Where("name = ?", ingredient.Name).FirstOrCreate(&ingredient)
	config.DB.Model(&recipe).Association("Ingredients").Append(&ingredient)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Adding new ingredient",
		"data":    ingredient,
	})
}

// DeleteIngredientFromRecipe
