package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrorRecordNotFound error. We'll return this from our
// Get() method then looking up a food doesn't exist in our database
var (
	ErrRecordNotFound = errors.New("Record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create models struct which wraps the FoodModel. We'll add another models to this
// like UserModel or PermissionModel, as our build progresses.
type Models struct {
	Foods FoodModel
	Users UserModel
	Token TokenModel
}

// For ease of use, we also add a New() method which return a Models struct constaining
// the initialized FoodModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Foods: FoodModel{DB: db},
		Users: UserModel{DB: db},
		Token: TokenModel{DB: db},
	}
}

// func NewModelsMock() Models {
// 	return Models{
// 		Foods: MockFoodModel{},
// 	}
// }
