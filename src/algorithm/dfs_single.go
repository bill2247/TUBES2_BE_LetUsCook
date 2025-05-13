package algorithm

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
)

func DFSSingle(t *data_type.RecipeTree) {
	currentId := scrapping.MapperNameToIdx[t.Name]

	// basis
	if currentId == 0 || currentId == 1 || currentId == 2 || currentId == 3 {
		return
	}

	// rekurens
	childrenList := scrapping.MapperIdxToRecipes[currentId]
	for i := 0; i < len(childrenList); i++ {
		idFirst := childrenList[i].First
		idSecond := childrenList[i].Second

		// skip yang tidak memenuhi kriteria
		if getTier(idFirst) >= getTier(currentId) || getTier(idSecond) >= getTier(currentId) {
			continue
		}

		firstRecipe := CreateRecipeTreeFromId(idFirst)
		secondRecipe := CreateRecipeTreeFromId(idSecond)

		currentPair := &data_type.Pair_recipe{First: firstRecipe, Second: secondRecipe}
		t.Children = append(t.Children, currentPair)
		DFSSingle(firstRecipe)
		DFSSingle(secondRecipe)
		break
	}
}

func DFSSingleEntryPoint(url string) (*data_type.RecipeTree, int) {
	idx, ok := scrapping.MapperNameToIdx[url]
	if idx == -1 || !ok{
		fmt.Println("Error: Invalid URL")
		return nil, 0
	}
	tier := scrapping.MapperIdxToTier[idx]
	if tier == -1 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx], Children: nil}, 1
	}
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx]}

	DFSSingle(root)
	return root, NodeCount(root)
}
