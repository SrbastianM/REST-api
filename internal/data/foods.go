package data

import "time"

type Food struct {
	ID       int64     // Unique integer for the Food
	CreateAt time.Time // TimeStamp when the Food is added to our db
	Title    string    // Food Title
	Types    []string  // Slices of types of food (Fruit and vegetables, starchy food, Dairy. Protein, fat)
	Version  int32     // Version number starts wiht 1 and will be incremented each time food information is updated
}
