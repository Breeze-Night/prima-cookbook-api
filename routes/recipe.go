package routes

import (
	"fmt"
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetRecipes(c *gin.Context) {
	recipes := []models.Recipe{}
	config.DB.Find(&recipes)

	config.DB.Preload(clause.Associations).Find(&recipes)

	c.JSON(http.StatusOK, gin.H{
		"message": "Find yourself an interesting recipe to try",
		"data":    recipes,
	})
}

func GetRecipeByID(c *gin.Context) {
	id := c.Param("id")

	var recipe models.Recipe

	data := config.DB.Preload(clause.Associations).First(&recipe, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Data Not Found",
			"message": "Recipe does not exist",
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Successful",
		"data":    recipe,
	})
}

func CreateRecipe(c *gin.Context) {
	var reqRecipe models.RecipeInput
	err := c.BindJSON(&reqRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad request",
		})
		c.Abort()
		return
	}

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

	recipe := models.Recipe{
		Title:        reqRecipe.Title,
		Description:  reqRecipe.Description,
		Instructions: reqRecipe.Instructions,
		UserID:       user.ID,
		Username:     user.Username,
	}

	if err := config.DB.Create(&recipe).Error; err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	// config.DB.Create(&recipe)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Success",
		"message": "Posting Recipe",
		"data":    recipe,
	})
}

// func EditRecipe(c *gin.Context) {
// 	recipeID := c.Param("id")

// 	var reqRecipe models.RecipeInput
// 	err := c.BindJSON(&reqRecipe)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   err.Error(),
// 			"message": "Bad request",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	var recipe models.Recipe
// 	err = config.DB.First(&recipe, recipeID).Error
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error":   "Recipe not found",
// 			"message": "Recipe not found",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	recipe.Title = reqRecipe.Title
// 	recipe.Description = reqRecipe.Description
// 	recipe.Instructions = reqRecipe.Instructions

// 	if err := config.DB.Save(&recipe).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"error":   err.Error(),
// 			"message": "Internal server error",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "Success",
// 		"message": "Updating Recipe",
// 		"data":    recipe,
// 	})
// }

func EditRecipe(c *gin.Context) {
	var reqRecipe models.RecipeInput
	err := c.BindJSON(&reqRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad request",
		})
		c.Abort()
		return
	}

	// Get recipe_id from path parameter
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid recipe id",
		})
		c.Abort()
		return
	}

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

	// Get recipe to update
	var recipe models.Recipe
	queryRes = config.DB.First(&recipe, recipeID)
	if queryRes.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Recipe with id %d not found", recipeID),
			"data":    "data not found",
		})
		return
	}

	// Check if the user is authorized to update the recipe
	if recipe.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "You are not authorized to update this recipe",
			"data":    "data not found",
		})
		return
	}

	recipe.Title = reqRecipe.Title
	recipe.Description = reqRecipe.Description
	recipe.Instructions = reqRecipe.Instructions
	recipe.Username = user.Username

	// Save updated recipe
	if err := config.DB.Save(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Recipe updated",
		"data":    recipe,
	})
}

// func DeleteRecipe(c *gin.Context) {
// 	recipeID := c.Param("id")

// 	var recipe models.Recipe
// 	err := config.DB.First(&recipe, recipeID).Error
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error":   "Recipe not found",
// 			"message": "Recipe not found",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	if err := config.DB.Delete(&recipe).Error; err != nil {
// 		c.JSON(500, gin.H{
// 			"error":   err.Error(),
// 			"message": "Internal server error",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "Success",
// 		"message": "Deleting Recipe",
// 		"data":    recipe,
// 	})
// }

func DeleteRecipe(c *gin.Context) {
	recipeID := c.Param("id")

	var recipe models.Recipe
	queryRes := config.DB.First(&recipe, "id = ?", recipeID)

	if queryRes.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Recipe with ID %s, not found", recipeID),
			"data":    "data not found",
		})
		return
	}

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
	queryRes = config.DB.First(&user, "email = ?", emailUser)

	if queryRes.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User by email %s, not found", emailUser),
			"data":    "data not found",
		})
		return
	}

	if recipe.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authorized to delete this recipe",
			"data":    "data not found",
		})
		return
	}

	if err := config.DB.Delete(&recipe).Error; err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Deleting Recipe",
		"data":    recipe,
	})
}
