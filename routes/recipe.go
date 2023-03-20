package routes

import (
	"net/http"
	"prima_cookbook/config"
	"prima_cookbook/models"

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

	recipe := models.Recipe{
		Title:        reqRecipe.Title,
		Description:  reqRecipe.Description,
		Instructions: reqRecipe.Instructions,
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
// 	id := c.Param("id")

// 	var reqRecipe models.RecipeInput
// 	c.BindJSON(&reqRecipe)

// 	recipe := models.Recipe{
// 		Title:        reqRecipe.Title,
// 		Description:  reqRecipe.Description,
// 		Instructions: reqRecipe.Instructions,
// 	}

// 	config.DB.Model(&recipe).Where("id = ?", id).Updates(reqRecipe)

//		c.JSON(200, gin.H{
//			"Message": "Recipe Updated",
//			"data":    recipe,
//		})
//	}

func EditRecipe(c *gin.Context) {
	recipeID := c.Param("id")

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

	var recipe models.Recipe
	err = config.DB.First(&recipe, recipeID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Recipe not found",
			"message": "Recipe not found",
		})
		c.Abort()
		return
	}

	recipe.Title = reqRecipe.Title
	recipe.Description = reqRecipe.Description
	recipe.Instructions = reqRecipe.Instructions

	if err := config.DB.Save(&recipe).Error; err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Updating Recipe",
		"data":    recipe,
	})
}

// func DeleteRecipe(c *gin.Context) {
// 	id := c.Param("id")
// 	var recipe models.Recipe

// 	data := config.DB.First(&recipe, "id = ?", id)

// 	if data.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"satus":   "Data Not Found",
// 			"message": "The recipe does not exist",
// 		})

// 		return
// 	}

// 	config.DB.Delete(&recipe, id)

// 	c.JSON(200, gin.H{
// 		"Message": "Recipe Deleted",
// 	})
// }

func DeleteRecipe(c *gin.Context) {
	recipeID := c.Param("id")

	var recipe models.Recipe
	err := config.DB.First(&recipe, recipeID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Recipe not found",
			"message": "Recipe not found",
		})
		c.Abort()
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
