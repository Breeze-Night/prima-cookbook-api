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

	// clean response output
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

	// clean response output
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

	// Get user_id in context by user email
	emailUser, exists := c.Get("x-email")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "x-email key not found in context",
			"message": "Bad request",
		})
		c.Abort()
		return
	}

	var user models.User
	queryRes := config.DB.Preload(clause.Associations).First(&user, "email = ?", emailUser)

	if queryRes.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User by email %s, not found", emailUser),
			"data":    "data not found",
		})
		return
	}

	// Check if the user is the owner of the recipe
	if recipe.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not authorized to perform this action",
		})
		return
	}

	err := c.BindJSON(&ingredient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})
		c.Abort()
		return
	}

	// check if ingredient already exist in database
	config.DB.Where("name = ?", ingredient.Name).FirstOrCreate(&ingredient)

	// adding the ingredient to the recipe's ingredients association
	config.DB.Model(&recipe).Association("Ingredients").Append(&ingredient)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Adding new ingredient",
		"data":    ingredient,
	})
}

func DeleteIngredientFromRecipe(c *gin.Context) {
	var recipe models.Recipe
	var ingredient models.Ingredient
	recipeID := c.Param("recipe_id")
	ingredientID := c.Param("ingredient_id")

	config.DB.First(&recipe, recipeID)

	// Get user_id in context by user email
	emailUser, exists := c.Get("x-email")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "x-email key not found in context",
			"message": "Bad request",
		})
		c.Abort()
		return
	}

	var user models.User
	queryRes := config.DB.Preload(clause.Associations).First(&user, "email = ?", emailUser)

	if queryRes.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User by email %s, not found", emailUser),
			"data":    "data not found",
		})
		return
	}

	// Check if the user is the owner of the recipe
	if recipe.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "You are not authorized to perform this action",
		})
		return
	}

	config.DB.First(&ingredient, ingredientID)

	if ingredient.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Ingredient not found",
		})
		return
	}

	// delete ingredient from the recipe by removing the association between the recipe and the ingredient
	config.DB.Model(&recipe).Association("Ingredients").Delete(&ingredient)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ingredient deleted",
	})
}
