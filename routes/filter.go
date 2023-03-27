package routes

import (
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"

	"github.com/gin-gonic/gin"
)

func GetRecipesByIngredient(c *gin.Context) {
	var recipes []models.Recipe
	ingredientID := c.Param("ingredient_id")

	config.DB.Where("id IN (SELECT recipe_id FROM recipe_ingredients WHERE ingredient_id = ?)", ingredientID).Find(&recipes)

	responseGetRecipe := []models.OutputAllRecipes{}

	for _, r := range recipes {
		aro := models.OutputAllRecipes{
			ID:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			Username:    r.Username,
		}
		responseGetRecipe = append(responseGetRecipe, aro)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseGetRecipe,
	})
}
