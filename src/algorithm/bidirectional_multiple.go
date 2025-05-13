package algorithm

import(
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
	"fmt"
)

func BidirectionalMultiple(targetURL string, bound int) (*data_type.RecipeTree, int) {
	targetIdx, ok := scrapping.MapperNameToIdx[targetURL]
	if !ok || targetIdx == -1 {
		fmt.Println("Error: Invalid target URL")
		return nil, 0
	}
	
	tier := scrapping.MapperIdxToTier[targetIdx]
	if tier == -1 || tier == 0 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[targetIdx], Children: nil}, 1
	}
	
	countRecipe := 0
	var allRecipes = [][]Path_element{}
	ForwardQueue := [][]Path_element{}
	BackwardQueue := [][]Path_element{}
	ForwardVisited := make(map[int][]Path_element) // ForwardVisited[x] menyimpan path yang menghasilkan x dari elemen dasar (ForwardVisited[x][].last() != x karena akan menghasilkan x)
	BackwardVisited := make(map[int][]Path_element) // BackwardVisited[(x,y)] menyimpan path yang dimulai dari pair (x,y) hingga menghasilkan targetIdx
	
	recipeMap:= scrapping.MapperPairToIdxs
	reverseMap := scrapping.MapperIdxToRecipes
	nodeCount := 0

	for i:=0; i<4; i++{
		if recipeMap[i] != nil {
			for key, val := range recipeMap[i] {
				result := val
				r1 := Path_element{Idx: i, Result: i}
				r2 := Path_element{Idx: key, Result: result}
				ForwardQueue = append(ForwardQueue, []Path_element{r1, r2})
				ForwardVisited[i] = []Path_element{r1}
				ForwardVisited[key] = []Path_element{r2}
				ForwardVisited[result] = []Path_element{r1, r2}
				nodeCount += 2
			}
		}
	}

	targetElmt := Path_element{Idx: targetIdx, Result: targetIdx}
	BackwardQueue = append(BackwardQueue, []Path_element{targetElmt})
	nodeCount += 1


	for len(ForwardQueue) > 0 && len(BackwardQueue) > 0 {
		// Proses Forward
		ForwardPath := ForwardQueue[0]
		ForwardQueue = ForwardQueue[1:]
		fwdlast := ForwardPath[len(ForwardPath)-1]
		fwdlastRes := fwdlast.Result
		nextElmts := recipeMap[fwdlastRes]
		for key, res := range nextElmts {
			if _, visited := ForwardVisited[res]; visited {
				continue 
			}
			nextPathElmt := Path_element{Idx: key, Result: res}
			tempForwardPath := append(ForwardPath, nextPathElmt)
			ForwardQueue = append(ForwardQueue, tempForwardPath)
			ForwardVisited[res] = tempForwardPath
			nodeCount ++
			if BackwardVisited[res] != nil {
				finalPath := []Path_element{}
				finalPath = append(finalPath, ForwardPath...)
				finalPath = append(finalPath, reversePath(BackwardVisited[res])...)
				allRecipes = append(allRecipes, finalPath)
				countRecipe ++
				if (countRecipe >= bound){
					break
				}
			}
		}

		// Backward process
		BackwardPath := BackwardQueue[0]
		BackwardQueue = BackwardQueue[1:]
		bwdlast := BackwardPath[len(BackwardPath)-1]
		bwdlastIdx := bwdlast.Idx
		recipes := reverseMap[bwdlastIdx]
		for _, recipe := range recipes {
	    	if _, visited1 := BackwardVisited[recipe.First]; visited1 {
				continue
			}
			if _, visited2 := BackwardVisited[recipe.Second]; visited2 {
				continue
    		}
			firstPathElmt := Path_element{Idx: recipe.First, Result: bwdlastIdx}
			secondPathElmt := Path_element{Idx: recipe.Second, Result: bwdlastIdx}
			firstPathElmtBasic := Path_element{Idx: recipe.First, Result: recipe.First}
			secondPathElmtBasic := Path_element{Idx: recipe.Second, Result: recipe.Second}
			newBackwardPath := BackwardPath[:len(BackwardPath)-1]
			newBackwardPath1 := append(newBackwardPath, firstPathElmt, secondPathElmtBasic)
			newBackwardPath2 := append(newBackwardPath, secondPathElmt, firstPathElmtBasic)
			BackwardQueue = append(BackwardQueue, newBackwardPath1)
			BackwardQueue = append(BackwardQueue, newBackwardPath2)
			BackwardVisited[recipe.Second] = append(BackwardVisited[recipe.Second], newBackwardPath1...)
			BackwardVisited[recipe.First] = append(BackwardVisited[recipe.First], newBackwardPath2...)
			nodeCount += 2
			if ForwardVisited[recipe.First] != nil {
				finalPath := []Path_element{}
				finalPath = append(finalPath, ForwardVisited[recipe.First]...)
				finalPath = append(finalPath, newBackwardPath2...)
				allRecipes = append(allRecipes, finalPath)
				countRecipe ++
				if (countRecipe >= bound){
					break
				}
			}
			if ForwardVisited[recipe.Second] != nil {
				finalPath := []Path_element{}
				finalPath = append(finalPath, ForwardVisited[recipe.Second]...)
				finalPath = append(finalPath, newBackwardPath1...)
				allRecipes = append(allRecipes, finalPath)
				countRecipe ++
				if (countRecipe >= bound){
					break
				}
			}
		}
	}
	
	if (len(allRecipes) == 0) {
		return nil, 0
	}
	baseRoot := pathTotree(allRecipes[0])
	for i := 1; i < len(allRecipes); i++ {
		combinePathToRoot(baseRoot, allRecipes[i], false)
	}
	baseRoot, newNodeCount := completeTheRoot(baseRoot)
	nodeCount += newNodeCount
	return baseRoot, nodeCount
}	

func combinePathToRoot(root *data_type.RecipeTree, recipe []Path_element, done bool) {
	if (done){
		return
	}
	if (root.Name != scrapping.MapperIdxToName[recipe[0].Idx]) {
		for _, child := range root.Children {
			combinePathToRoot(child.First, recipe, done)
		}
	} else {
		newRecipe := pathTotree(recipe)
		root.Children = append(root.Children, newRecipe.Children...)
	}
}