package dfs

// this is an unused, legacy file

// var visited = make([]bool, 721)
// var MaxDepth = 999999

// func SetMaxDepth(x int) {
// 	MaxDepth = x
// }

// func CreateRecipeTreeFromName(name string) *data_type.RecipeTreeWithDepth {
// 	rootChild := []*data_type.PairRecipeWithDepth{}
// 	visitedMap := make(map[int]bool)
// 	return &data_type.RecipeTreeWithDepth{Name: name, Children: rootChild, Depth: -1, NextHop: -1, Visited: visitedMap}
// }

// func CreateRecipeTreeFromId(id int) *data_type.RecipeTreeWithDepth {
// 	rootChild := []*data_type.PairRecipeWithDepth{}
// 	name := scrapping.MapperElmIdx[id]
// 	visitedMap := make(map[int]bool)
// 	return &data_type.RecipeTreeWithDepth{Name: name, Children: rootChild, Depth: -1, NextHop: -1, Visited: visitedMap}
// }

// func setVisited(t *data_type.RecipeTreeWithDepth, v map[int]bool) {
// 	for key, value := range v {
// 		t.Visited[key] = value
// 	}
// }

// func max2(a int, b int) int {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// }

// func min2(a int, b int) int {
// 	if a < b {
// 		return a
// 	} else {
// 		return b
// 	}
// }

// func PrintShortestPath(t *data_type.RecipeTreeWithDepth, indentCount int) {
// 	fmt.Print(strings.Repeat(" ", 2*indentCount))
// 	fmt.Println(t.Name + " (Depth: " + strconv.Itoa(t.Depth) + " Next Hop: " + strconv.Itoa(t.NextHop) + ")")
// 	if t.NextHop == -1 {
// 		return
// 	}
// 	PrintShortestPath(t.Children[t.NextHop].First, indentCount+1)
// 	PrintShortestPath(t.Children[t.NextHop].Second, indentCount+1)
// }

// func DFSShortestPath(t *data_type.RecipeTreeWithDepth) {
// 	// fmt.Print(" ")    // Debugging purposes
// 	fmt.Print(t.Name) // Debugging purposes
// 	// reader := bufio.NewReader(os.Stdin)
// 	// input, _ := reader.ReadString('\n')
// 	// input += "a"
// 	code := scrapping.MapperIdxElm[t.Name]

// 	// basis
// 	if code == 0 || code == 1 || code == 2 || code == 3 {
// 		t.Depth = 0
// 		return
// 	}
// 	// handle rekursif
// 	if visited, ok := t.Visited[code]; ok && visited {
// 		return
// 	}

// 	// rekurens
// 	t.Visited[code] = true
// 	childrenList := scrapping.MapperRecipe1[code]

// 	// preserve the old map
// 	legacyVisited := make(map[int]bool)
// 	for key, value := range t.Visited {
// 		legacyVisited[key] = value
// 	}

// 	for i := 0; i < len(childrenList); i++ {
// 		// fmt.Print(" " + t.Name)      // debugging purposes
// 		setVisited(t, legacyVisited) // reset the visited map

// 		firstRecipe := CreateRecipeTreeFromId(childrenList[i].First)
// 		setVisited(firstRecipe, t.Visited)
// 		secondRecipe := CreateRecipeTreeFromId(childrenList[i].Second)
// 		setVisited(secondRecipe, t.Visited)

// 		currentPair := &data_type.PairRecipeWithDepth{First: firstRecipe, Second: secondRecipe}
// 		t.Children = append(t.Children, currentPair)

// 		// boleh pruning? kalau depthnya udah lebih, stop (nanti akalin pake global aja)
// 		DFSShortestPath(firstRecipe)
// 		DFSShortestPath(secondRecipe)

// 		if firstRecipe.Depth != -1 && secondRecipe.Depth != -1 {
// 			childDepth := max2(firstRecipe.Depth, secondRecipe.Depth)
// 			if t.Depth == -1 {
// 				t.Depth = childDepth + 1
// 				t.NextHop = i
// 			} else {
// 				if childDepth < t.Depth && childDepth != -1 {
// 					t.NextHop = i
// 					t.Depth = min2(t.Depth, childDepth)
// 				}
// 			}
// 		}
// 	}
// }
