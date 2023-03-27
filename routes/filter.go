package routes

import (
	"fmt"
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"

	"github.com/gin-gonic/gin"
)

func GetRecipesByIngredient(c *gin.Context) {
	var recipes []models.Recipe
	ingredientID := c.Param("ingredient_id")

	// use subquery to retrieve all recipeID from recipe_ingredients table where the ingredient_id matches the provided value
	// then retrieving the full recipe details from the recipes table where the ID is in the subquery
	config.DB.Where("id IN (SELECT recipe_id FROM recipe_ingredients WHERE ingredient_id = ?)", ingredientID).Find(&recipes)

	// response diperpendek untuk enak dibaca
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
		"message": fmt.Sprintf("Here are all recipes with ingredient_id: %s", ingredientID),
		"data":    responseGetRecipe,
	})
}

func GetRecipeByUserID(c *gin.Context) {
	var recipes []models.Recipe
	userID := c.Param("user_id")

	// retrieve all recipes associated with the specified user ID
	config.DB.Where("user_id = ?", userID).Find(&recipes)

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
