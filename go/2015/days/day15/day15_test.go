package main

import (
	"reflect"
	"testing"
)

// the exemple in the problem text
var stock Stock = Stock{
	map[string]int{"Butterscotch": 0, "Cinnamon": 1},
	[]string{"Butterscotch", "Cinnamon"},
	[]Ingredient{
		Ingredient{"Butterscotch", 0, -1, -2, 6, 3, 8},
		Ingredient{"Cinnamon", 1, 2, 3, -2, -1, 3},
	},
}

// reflect.DeepEqual(got, expected)
func Test_parseIngredients(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := parseIngredients([]string{("Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3")})
		expected := &Stock{
			map[string]int{"Cinnamon": 0},
			[]string{"Cinnamon"},
			[]Ingredient{Ingredient{"Cinnamon", 0, 2, 3, -2, -1, 3}},
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := parseIngredients(fileToLines("exemple.txt"))
		expected := &stock
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_recipeScore(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := recipeScore(&stock, []int{44, 56})
		expected := 62842880
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := recipeScore(&stock, []int{1, 99})
		expected := 0
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
