package bfs_multiple_recipe

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
	"strings"
	"sync"
	"let_us_cook/src/bfs_shortest"
)

type bfsTask struct {
	Idx   int
	Node  *data_type.RecipeTree
	Depth int
}

func Bfs_multiple_recipe(url string, bound int) (*data_type.RecipeTree, int) {
	idx := scrapping.MapperNameToIdx[url]
	if idx == -1 {
		fmt.Println("Error: Invalid URL")
		return nil, 0
	}
	tier := scrapping.MapperIdxToTier[idx]
	if tier == -1 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx], Children: nil}, 1
	}

	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx]}
	queue := make(chan bfsTask, 1000)
	visited := make([]bool, 720)
	var visitedMu sync.Mutex

	countNode := 1
	var count int
	var countMu sync.Mutex

	var wg sync.WaitGroup
	queue <- bfsTask{Idx: idx, Node: root, Depth: 0}
	wg.Add(1)

	for i := 0; i < 8; i++ { // 8 workers
		go func() {
			for task := range queue {
				idx := task.Idx
				node := task.Node
				depth := task.Depth

				visitedMu.Lock()
				if visited[idx] {
					visitedMu.Unlock()
					wg.Done()
					continue
				}
				visited[idx] = true
				visitedMu.Unlock()

				if idx <= 4 {
					wg.Done()
					continue
				}

				recipes := scrapping.MapperIdxToRecipes[idx]
				for _, recipe := range recipes {
					firstIdx := recipe.First
					secondIdx := recipe.Second

					if scrapping.MapperIdxToTier[firstIdx] >= tier || scrapping.MapperIdxToTier[secondIdx] >= tier {
						continue
					}

					countMu.Lock()
					if count >= bound {
						countMu.Unlock()
						newRoot, nodeAddition := completeTheRoot(root)
						root = newRoot
						countNode += nodeAddition
						break
					}
					if countRecipe(root) >= bound {
						count++
					}
					countMu.Unlock()
					// itung count node
					countNode += 2
					firstNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[firstIdx]}
					secondNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[secondIdx]}
					pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}
					node.Children = append(node.Children, pair)

					wg.Add(2)
					queue <- bfsTask{Idx: firstIdx, Node: firstNode, Depth: depth + 1}
					queue <- bfsTask{Idx: secondIdx, Node: secondNode, Depth: depth + 1}
				}
				wg.Done()
			}
		}()
	}

	wg.Wait()
	close(queue)

	// Prune setelah BFS selesai
	PruneNonTerminalParallel(root)
	PruneNonTerminalParallel(root)

	return root, countNode // debugging
}

// ------------------------- PARALLEL PRUNE ----------------------------

func PruneNonTerminalParallel(root *data_type.RecipeTree) {
	pruneTreeParallel(root)
}

func pruneTreeParallel(node *data_type.RecipeTree) bool {
	if node == nil {
		return false
	}

	if len(node.Children) == 0 {
		return isBasicElement(node)
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	validChildren := make([]*data_type.Pair_recipe, 0)

	for _, pair := range node.Children {
		wg.Add(1)
		go func(pair *data_type.Pair_recipe) {
			defer wg.Done()
			firstValid := pruneTreeParallel(pair.First)
			secondValid := pruneTreeParallel(pair.Second)

			if len(pair.First.Children) == 0 && len(pair.Second.Children) == 0 {
				if isBasicElement(pair.First) && isBasicElement(pair.Second) {
					mu.Lock()
					validChildren = append(validChildren, pair)
					mu.Unlock()
				}
			} else if firstValid || secondValid {
				mu.Lock()
				validChildren = append(validChildren, pair)
				mu.Unlock()
			}
		}(pair)
	}

	wg.Wait()
	node.Children = validChildren
	return len(validChildren) > 0
}

func isBasicElement(node *data_type.RecipeTree) bool {
	idx, ok := scrapping.MapperNameToIdx[node.Name]
	return ok && idx <= 4
}

func TreeToString(node *data_type.RecipeTree) string {
	var builder strings.Builder
	displayTreeToBuilder(node, "", true, &builder)
	return builder.String()
}

func displayTreeToBuilder(node *data_type.RecipeTree, prefix string, isTail bool, builder *strings.Builder) {
	if node == nil {
		return
	}
	builder.WriteString(prefix + branchSymbol(isTail) + node.Name + "\n")

	children := node.Children
	for i, pair := range children {
		isLast := i == len(children)-1
		displayTreeToBuilder(pair.First, prefix+nextPrefix(isTail), false, builder)
		displayTreeToBuilder(pair.Second, prefix+nextPrefix(isTail), isLast, builder)
	}
}

// DisplayTree dan fungsi pendukung tetap sama
func DisplayTree(node *data_type.RecipeTree, prefix string, isTail bool) {
	if node == nil {
		return
	}
	fmt.Println(prefix + branchSymbol(isTail) + node.Name)

	children := node.Children
	for i, pair := range children {
		isLast := i == len(children)-1
		DisplayTree(pair.First, prefix+nextPrefix(isTail), false)
		DisplayTree(pair.Second, prefix+nextPrefix(isTail), isLast)
	}
}

func branchSymbol(isTail bool) string {
	if isTail {
		return "└── "
	}
	return "├── "
}

func nextPrefix(isTail bool) string {
	if isTail {
		return "    "
	}
	return "│   "
}

func countRecipe(t *data_type.RecipeTree) int {
	currentId := scrapping.MapperNameToIdx[t.Name]

	if currentId == 0 || currentId == 1 || currentId == 2 || currentId == 3 {
		return 1
	}

	// rekurens
	totalWays := 0
	for i := 0; i < len(t.Children); i++ {
		firstRecipe := t.Children[i].First
		secondRecipe := t.Children[i].Second

		countFirst := countRecipe(firstRecipe)
		countSecond := countRecipe(secondRecipe)

		totalWays += countFirst * countSecond
	}
	return totalWays
}

func completeTheRoot(root *data_type.RecipeTree) (*data_type.RecipeTree, int) {
	idx := scrapping.MapperNameToIdx[root.Name]
	if idx > 4 || len(root.Children) == 0 {
		newRoot, nodeCount := bfs_shortest.FindShortestPath(root.Name)
		return newRoot, nodeCount
	} else {
		nodeCount := 0
		for i := 0; i < len(root.Children); i++ {
			firstRecipe, firstNodeCount := completeTheRoot(root.Children[i].First)
			secondRecipe, secondNodeCount := completeTheRoot(root.Children[i].Second)

			root.Children[i].First = firstRecipe
			root.Children[i].Second = secondRecipe
			nodeCount += firstNodeCount + secondNodeCount
		}
		return root, nodeCount
	}
}