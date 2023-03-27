package routes

import (
	"fmt"
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetIngredients(c *gin.Context) {
	ingredients := []models.Ingredient{}

	config.DB.Preload(clause.Associations).Find(&ingredients)

	responseGetIngredients := []models.OutputIngredients{}

	for _, ingredient := range ingredients {
		responseRecipes := []models.RecipesInIngredient{}
		for _, recipe := range ingredient.Recipes {
			rii := models.RecipesInIngredient{
				ID:    recipe.ID,
				Title: recipe.Title,
			}
			responseRecipes = append(responseRecipes, rii)
		}

		oai := models.OutputIngredients{
			ID:      ingredient.ID,
			Name:    ingredient.Name,
			Recipes: responseRecipes,
		}
		responseGetIngredients = append(responseGetIngredients, oai)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All registered ingredients",
		"data":    responseGetIngredients,
	})
}

func GetIngredientByID(c *gin.Context) {
	id := c.Param("id")

	var ingredient models.Ingredient

	data := config.DB.Preload(clause.Associations).First(&ingredient, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Data Not Found",
			"message": fmt.Sprintf("Ingredient with id: %s does not exist", id),
		})
		return
	}

	responseRecipes := []models.RecipesInIngredient{}
	for _, recipe := range ingredient.Recipes {
		rii := models.RecipesInIngredient{
			ID:    recipe.ID,
			Title: recipe.Title,
		}
		responseRecipes = append(responseRecipes, rii)
	}

	responseGetIngredients := models.OutputIngredients{
		ID:      ingredient.ID,
		Name:    ingredient.Name,
		Recipes: responseRecipes,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
		"data":    responseGetIngredients,
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
