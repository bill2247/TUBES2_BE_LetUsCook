package dfs

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scrapping"
	"strings"
)

func CreateRecipeTreeFromName(name string) *data_type.RecipeTree {
	rootChild := []*data_type.Pair_recipe{}
	return &data_type.RecipeTree{Name: name, Children: rootChild}
}

func CreateRecipeTreeFromId(id int) *data_type.RecipeTree {
	rootChild := []*data_type.Pair_recipe{}
	name := scrapping.MapperElmIdx[id]
	return &data_type.RecipeTree{Name: name, Children: rootChild}
}

func PrintTree(t *data_type.RecipeTree, indentCount int) {
	fmt.Print(strings.Repeat(" ", 2*indentCount))
	fmt.Println(t.Name)
	if indentCount == 1 {
		return
	}
	for i := 0; i < len(t.Children); i++ {
		PrintTree(t.Children[i].First, indentCount+1)
		PrintTree(t.Children[i].Second, indentCount+1)
	}

}

func NodeCount(t *data_type.RecipeTree) int {
	currentId := scrapping.MapperIdxElm[t.Name]
	if currentId == 0 || currentId == 1 || currentId == 2 || currentId == 3 || getTier(currentId) == 9999 {
		return 1
	}

	total := 1
	for i := 0; i < len(t.Children); i++ {
		total += NodeCount(t.Children[i].First)
		total += NodeCount(t.Children[i].Second)
	}
	return total
}

func getTier(id int) int {
	if id <= 3 { // water
		return 0
	} else if id == 4 { // time
		return 9999
	} else if id <= 14 { // steam
		return 1
	} else if id <= 27 { // wind
		return 2
	} else if id <= 39 { // warmth
		return 3
	} else if id <= 68 { // wave
		return 4
	} else if id <= 112 { // windmill
		return 5
	} else if id <= 148 { // water gun
		return 6
	} else if id <= 172 { // wire
		return 7
	} else if id <= 249 { // wolf
		return 8
	} else if id <= 384 { // zombie
		return 9
	} else if id <= 501 { // zoo
		return 10
	} else if id <= 603 { // yogurt
		return 11
	} else if id <= 675 { // writer
		return 12
	} else if id <= 709 { // vinegar
		return 13
	} else if id <= 717 { // treasure
		return 14
	} else if id <= 719 { // yogurt
		return 15
	}
	return 9999
}
