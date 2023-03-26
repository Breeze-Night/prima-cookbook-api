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
