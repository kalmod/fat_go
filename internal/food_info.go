package food

import "strings"

type Macro struct {
	Name   string
	Amount float64
	Metric Unit // I'll probably want to turn this into a type def?
}

type ServingSize struct {
	Amount int
	Metric Unit
}

type FoodNutrition struct {
	FoodName  string
	FoodID    string
	BrandName string
	Serving   ServingSize
	Macros    []Macro
}

type Unit string

const (
	grams   Unit = "g"
	oz      Unit = "oz"
	serving Unit = "serv"
	kcal    Unit = "kcal"
)

// if it's not grams or oz, it's probably a serving of itself.
// i.e 1 banana or 1 egg. So the type will just be serving.
func GetUnitType(unit string) Unit {

	switch unitType := strings.ToLower(unit); Unit(unitType) {
	case oz:
		return oz
	case grams:
		return grams
	case kcal:
		return kcal
	default:
		return serving
	}
}
