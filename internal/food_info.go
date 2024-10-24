package food

import (
	"fmt"
	"strings"
)

type Macro struct {
	Name   Macronutrient
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
type Macronutrient string

const (
	grams   Unit          = "g"
	oz      Unit          = "oz"
	serving Unit          = "serv"
	kcal    Unit          = "kcal"
	protein Macronutrient = "protein"
	carbs   Macronutrient = "carbs"
	fat     Macronutrient = "fat"
	cal     Macronutrient = "calories"
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

func GetMacroType(macro string) Macronutrient {
	switch macrotype := strings.ToLower(macro); Macronutrient(macrotype) {
	case protein:
		return protein
	case carbs:
		return carbs
	case fat:
		return fat
	case cal:
		return cal
	default:
		return Macronutrient("")
	}
}

func (m Macronutrient) MacroStringColor() string {
	switch m {
	case protein:
		return fmt.Sprintf("\x1b[33m%v\x1b[0m", m)
	case carbs:
		return fmt.Sprintf("\x1b[32m%v\x1b[0m", m)
	case fat:
		return fmt.Sprintf("\x1b[31m%v\x1b[0m", m)
	case cal:
		return fmt.Sprintf("\x1b[36m%v\x1b[0m", m)
	default:
		return fmt.Sprintf("%v", m)
	}
}

func (m *Macro) PrettyPrintMacro() string {
	return fmt.Sprintf("%v: %.2f%v", m.Name.MacroStringColor(), m.Amount, m.Metric)
}

func (fn *FoodNutrition) PrettyPrintNutrition() string {
	var prettyString string
	prettyString = prettyString + fmt.Sprintf("\x1b[3m%v\x1b[0m (%v) -- %v %v\n", fn.FoodName, fn.FoodID, fn.Serving.Amount, fn.Serving.Metric)
	for _, mc := range fn.Macros {
		prettyString = prettyString + mc.PrettyPrintMacro() + " "
	}
	return prettyString
}

func GetNewNutritionByWeight(old FoodNutrition, new_amount int) FoodNutrition {
	newServingSize := ServingSize{Amount: (old.Serving.Amount / old.Serving.Amount) * new_amount, Metric: old.Serving.Metric}
	newMacros := []Macro{}
	for _, mc := range old.Macros {
		newMacros = append(newMacros, Macro{mc.Name, (mc.Amount / float64(old.Serving.Amount)) * float64(new_amount), mc.Metric})
	}

	return FoodNutrition{FoodName: old.FoodName, FoodID: old.FoodID, BrandName: old.BrandName, Serving: newServingSize, Macros: newMacros}
}
