package models

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model
	Name    string   `json:"name"`
	Recipes []Recipe `gorm:"many2many:recipe_ingredients;" json:"recipes"`
}

type OutputIngredients struct {
	ID      uint                  `json:"id"`
	Name    string                `json:"name"`
	Recipes []RecipesInIngredient `json:"recipes"`
}

type RecipesInIngredient struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
