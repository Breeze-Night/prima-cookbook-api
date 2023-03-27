package models

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Instructions string       `json:"instructions"`
	Ingredients  []Ingredient `gorm:"many2many:recipe_ingredients;" json:"ingredients"`
	UserID       uint         `json:"user_id"`
	Username     string       `json:"username"`
}

type RecipeInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

type OutputAllRecipes struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Username    string `json:"username"`
}

type OutputRecipeByID struct {
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	Username     string                `json:"username"`
	Instructions string                `json:"instructions"`
	Ingredients  []IngredientsInRecipe `json:"ingredients"`
}

type IngredientsInRecipe struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
