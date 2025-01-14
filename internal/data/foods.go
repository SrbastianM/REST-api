package data

import (
	"SrbastianM/rest-api-gin/internal/validator"
	"database/sql"
	"time"
)

type Food struct {
	ID       int64     // Unique integer for the Food
	CreateAt time.Time // TimeStamp when the Food is added to our db
	Title    string    // Food Title
	Types    []string  // Slices of types of food (Fruit and vegetables, starchy food, Dairy. Protein, fat)
	Version  int32     // Version number starts wiht 1 and will be incremented each time food information is updated
}

// Define a FoodModel Struct type which wraps a sql.DB connection pools.
type FoodModel struct {
	DB *sql.DB
}

// // Create mock to Unit test all of the methods: Create, Get, Update and Delete
// type MockFoodModel struct{}

// func (f MockFoodModel) Insert(food *Food) error {
// 	return nil
// }

// func (f MockFoodModel) Get(id int64) (*Food, error) {
// 	return nil, nil
// }

// func (f MockFoodModel) Update(food *Food) error {
// 	return nil
// }

// func (f MockFoodModel) Delete(id int64) error {
// 	return nil
// }

// Add placeholder method for inserting a new record in the food table.
func (f FoodModel) Insert(food *Food) error {
	query :=
		`INSERT INTO foods (title, type)
	 VALUES ($1, $2)
	 RETURNING id, created_at, version
	`
	args := []interface{}{food.Title, food.Types}
	return f.DB.QueryRow(query, args...).Scan(&food.ID, &food.CreateAt, &food.Version)
}

// Add placeholder method for fetching a specific record from the food table.
func (f FoodModel) Get(id int64) (*Food, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the food table.
func (f FoodModel) Update(food *Food) error {
	return nil
}

// Add a placeholder method for deleting a specific record from movies table.
func (f FoodModel) Delete(id int64) (*Food, error) {
	return nil, nil
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
