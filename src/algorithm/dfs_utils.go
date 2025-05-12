package algorithm

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
	"strings"
)

func CreateRecipeTreeFromName(name string) *data_type.RecipeTree {
	rootChild := []*data_type.Pair_recipe{}
	return &data_type.RecipeTree{Name: name, Children: rootChild}
}

func CreateRecipeTreeFromId(id int) *data_type.RecipeTree {
	rootChild := []*data_type.Pair_recipe{}
	name := scrapping.MapperIdxToName[id]
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
	currentId := scrapping.MapperNameToIdx[t.Name]
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
	return scrapping.MapperIdxToTier[id]
}
