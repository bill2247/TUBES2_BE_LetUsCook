package algorithm

import(
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
	"fmt"
)

type Path_element struct {
	Idx int
	Result []int
}

func BidirectionalSingle(targetURL string) (*data_type.RecipeTree, int) {
	targetIdx, ok := scrapping.MapperNameToIdx[targetURL]
	if !ok || targetIdx == -1 {
		fmt.Println("Error: Invalid target URL")
		return nil, 0
	}

	tier := scrapping.MapperIdxToTier[targetIdx]
	if tier == -1 || tier == 0 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[targetIdx], Children: nil}, 1
	}

	ForwardQueue := [][]Path_element{}
	BackwardQueue := [][]Path_element{}
	ForwardVisited := make(map[int][][]Path_element) // ForwardVisited[x] menyimpan path yang menghasilkan x dari elemen dasar (ForwardVisited[x][].last() != x karena akan menghasilkan x)
	BackwardVisited := make(map[data_type.Recipe][][]Path_element) // BackwardVisited[(x,y)] menyimpan path yang dimulai dari pair (x,y) hingga menghasilkan targetIdx
	
	nodeCount := 0

	for i:= 0; i <= 4; i++ {
		for j:=0; j <= 4; j++ {
			x1 := data_type.Recipe{First: i, Second: j}
			x2 := data_type.Recipe{First: j, Second: i}
			if scrapping.MapperPairToIdxs[x1] != nil || scrapping.MapperPairToIdxs[x2] != nil {
				result := scrapping.MapperPairToIdxs[data_type.Recipe{First: i, Second: j}]
				r1 := Path_element{Idx: i, Result: []int{i}}
				r2 := Path_element{Idx: j, Result: result}
				ForwardQueue = append(ForwardQueue,[]Path_element{r1, r2})
				ForwardVisited[i] = [][]Path_element{{r1}}
				ForwardVisited[j] = [][]Path_element{{r2}}
				for _, v := range result {
					ForwardVisited[v] = [][]Path_element{{r1, r2}}
				}
				nodeCount += 2
			}
		}
	}

	targetElmt := Path_element{Idx: targetIdx, Result: []int{targetIdx}}
	BackwardQueue = append(BackwardQueue, []Path_element{targetElmt})
	nodeCount += 1

	recipeMap:= scrapping.MapperPairToIdxs
	reverseMap := scrapping.MapperIdxToRecipes

	for len(ForwardQueue) > 0 && len(BackwardQueue) > 0 {
		// Proses Forward
		ForwardPath := ForwardQueue[0]
		ForwardQueue = ForwardQueue[1:]
		fwdlast := ForwardPath[len(ForwardPath)-1]
		fwdlastRes := fwdlast.Result

		// Traversal recipeMap juga traversal fwdlastRes. 
		// Jika fwdlastRes[i] ada di recipeMap, maka dilakukan:
		// 1. Edit Result dari Path_element fwdlast menjadi tunggal yaitu fwdlastRes[i]
		// 3. Copy ForwardPath ke tempForwardPath, hapus elemen terakhir,
		// 4. bangan nextPathElmt yaitu Path_element{Idx: k.Second, Result: res}
		// 5. Append tempForwardPath dengan newPathElmt dan nextPathElmt
		// 6. Append tempForwardPath ke ForwardQueue
		// 7. Untuk setiap elemen di res, tambahkan ForwardPath ke ForwardVisited[res[i]]

		for k, res := range recipeMap {
			for i := 0; i < len(fwdlastRes); i++ {
				if fwdlastRes[i] == k.First{
					newPathElmt := Path_element{Idx: fwdlast.Idx, Result: []int{fwdlastRes[i]}}
					nextPathElmt := Path_element{Idx: k.Second, Result: res}
					tempForwardPath := ForwardPath[:len(ForwardPath)-1] 
					tempForwardPath = append(tempForwardPath, newPathElmt, nextPathElmt)
					ForwardQueue = append(ForwardQueue, tempForwardPath)
					nodeCount ++
					for _, v := range res {
						ForwardVisited[v] = append(ForwardVisited[v], ForwardPath)
						for pair, paths := range BackwardVisited {
							if pair.First == v {
								FinalPaths := [][]Path_element{}
								for i := 0; i < len(paths); i++ {
									FinalPath := []Path_element{}
									FinalPath = append(FinalPath, ForwardPath...)
									FinalPath = append(FinalPath, reversePath(paths[i])...)
									FinalPaths = append(FinalPaths, FinalPath)
								}
								// cari path terpendek
								if len(FinalPaths) > 0 {
									minPath := FinalPaths[0]
									for _, p := range FinalPaths {
										if len(p) < len(minPath) {
											minPath = p
										}
									}
									// print path
									fmt.Println("Found path:")
									for _, elmt := range minPath {
										fmt.Printf("%d ", elmt.Idx)
									}
									fmt.Println()
									root := pathTotree(minPath)
									newNodeCount := 0
									root, newNodeCount = completeTheRoot(root)
									nodeCount += newNodeCount
									return root, nodeCount
								}
							}
						}
					}

				} else if fwdlastRes[i] == k.Second {
					newPathElmt := Path_element{Idx: fwdlast.Idx, Result: []int{fwdlastRes[i]}}
					nextPathElmt := Path_element{Idx: k.First, Result: res}
					tempForwardPath := ForwardPath[:len(ForwardPath)-1]
					tempForwardPath = append(tempForwardPath, newPathElmt, nextPathElmt)
					ForwardQueue = append(ForwardQueue, tempForwardPath)
					nodeCount ++
					for _, v := range res {
						ForwardVisited[v] = append(ForwardVisited[v], ForwardPath)
						for pair, paths := range BackwardVisited {
							if pair.Second == v {
								FinalPaths := [][]Path_element{}
								for i := 0; i < len(paths); i++ {
									FinalPath := []Path_element{}
									FinalPath = append(FinalPath, ForwardPath...)
									FinalPath = append(FinalPath, reversePath(paths[i])...)
									FinalPaths = append(FinalPaths, FinalPath)
								}
								// cari path terpendek
								if len(FinalPaths) > 0 {
									minPath := FinalPaths[0]
									for _, p := range FinalPaths {
										if len(p) < len(minPath) {
											minPath = p
										}
									}
									// print path
									fmt.Println("Found path:")
									for _, elmt := range minPath {
										fmt.Printf("%d ", elmt.Idx)
									}
									fmt.Println()
									root := pathTotree(minPath)
									newNodeCount := 0
									root, newNodeCount = completeTheRoot(root)
									nodeCount += newNodeCount
									return root, nodeCount
								}
							}
						}
					}
				} else {
					continue
				}
			}
		}

		// Proses Backward
		BackwardPath := BackwardQueue[0]
		BackwardQueue = BackwardQueue[1:]
		bwdlast := BackwardPath[len(BackwardPath)-1]
		bwdlastIdx := bwdlast.Idx

		recipes := reverseMap[bwdlastIdx]
		for _, recipe := range recipes {
			// recipe = (first, second)
			// Jika path saat ini adalah {X} - currentElmt
			// Update path menjadi {X} - fisrt - second
			firstPathElmt := Path_element{Idx: recipe.First, Result: []int{bwdlastIdx}}
			secondPathElmt := Path_element{Idx: recipe.Second, Result: []int{bwdlastIdx}}
			firstPathElmtBasic := Path_element{Idx: recipe.First, Result: []int{recipe.First}}
			secondPathElmtBasic := Path_element{Idx: recipe.Second, Result: []int{recipe.Second}}
			newBackwardPath := BackwardPath[:len(BackwardPath)-1]
			newBackwardPath = append(newBackwardPath, firstPathElmt, secondPathElmtBasic)
			newBackwardPath = append(newBackwardPath, secondPathElmt, firstPathElmtBasic)
			BackwardQueue = append(BackwardQueue, newBackwardPath)
			nodeCount += 2
			// Update BackwardVisited
			BackwardVisited[recipe] = append(BackwardVisited[recipe], newBackwardPath)
			// cek apakah ForwardVisited[recipe.First] atau ForwardVisited[recipe.Second] ada yang sama dengan BackwardVisited[recipe]
			FinalPaths := [][]Path_element{}
			for i := 0; i < len(ForwardVisited[recipe.First]); i++ {
				for j := 0; j < len(BackwardVisited[recipe]); j++ {
					FinalPath := []Path_element{}
					FinalPath = append(FinalPath, ForwardVisited[recipe.First][i]...)
					FinalPath = append(FinalPath, reversePath(BackwardVisited[recipe][j])...)
					FinalPaths = append(FinalPaths, FinalPath)
				}
			}
			for i := 0; i < len(ForwardVisited[recipe.Second]); i++ {
				for j := 0; j < len(BackwardVisited[recipe]); j++ {
					FinalPath := []Path_element{}
					FinalPath = append(FinalPath, ForwardVisited[recipe.Second][i]...)
					FinalPath = append(FinalPath, reversePath(BackwardVisited[recipe][j])...)
					FinalPaths = append(FinalPaths, FinalPath)
				}
			}
			// cari path terpendek
			if len(FinalPaths) > 0 {
				minPath := FinalPaths[0]
				for _, p := range FinalPaths {
					if len(p) < len(minPath) {
						minPath = p
					}
				}
				// print path
				fmt.Println("Found path:")
				for _, elmt := range minPath {
					fmt.Printf("%d ", elmt.Idx)
				}
				fmt.Println()
				root := pathTotree(minPath)
				newNodeCount := 0
				root, newNodeCount = completeTheRoot(root)
				nodeCount += newNodeCount
				return root, nodeCount
			}
		}
		
	}
	return nil, nodeCount
}	

func reversePath(path []Path_element) []Path_element {
	reversed := make([]Path_element, len(path))
	for i, j := 0, len(path)-1; i <= j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = path[j], path[i]
	}
	return reversed
}

func pathTotree(path []Path_element) *data_type.RecipeTree {
	if len(path) == 0 {
		return nil
	}
	path = reversePath(path)
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[path[0].Result[0]], Children: nil}
	prev := root
	for i, _ := range path[:len(path)-1] {
		firstRecipe := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[path[i].Idx], Children: nil}
		secondRecipe := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[path[i+1].Result[0]], Children: nil}
		newRecipe := &data_type.Pair_recipe{First: firstRecipe, Second: secondRecipe}
		prev.Children = append(prev.Children, newRecipe)
		prev = firstRecipe
	}
	return root
}