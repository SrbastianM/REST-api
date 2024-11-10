package data

import (
	"SrbastianM/rest-api-gin/internal/validator"
	"time"
)

type Food struct {
	ID       int64     // Unique integer for the Food
	CreateAt time.Time // TimeStamp when the Food is added to our db
	Title    string    // Food Title
	Types    []string  // Slices of types of food (Fruit and vegetables, starchy food, Dairy. Protein, fat)
	Version  int32     // Version number starts wiht 1 and will be incremented each time food information is updated
}

func ValidateFood(v *validator.Validator, food *Food) {
	//Use Check() method to execute the validation checks -> See the validator on internal/validator/validator.go
	v.Check(food.Title != "", "title", "must be provided")
	v.Check(len(food.Title) <= 500, "title", "must not be more than 500 bytes long")
	// Types
	v.Check(food.Types != nil, "types", "must be provided")
	v.Check(len(food.Types) >= 1, "types", "must be contained 1 type")
	v.Check(len(food.Types) <= 5, "types", "must not contain more than 5 type")

	v.Check(validator.Unique(food.Types), "types", "must not contain duplicate values")
}
