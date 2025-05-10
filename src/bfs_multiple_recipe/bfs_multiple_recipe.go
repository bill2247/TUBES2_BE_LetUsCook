package bfs_multiple_recipe

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scrapping"
	"strings"
)

type bfsTask struct {
	Idx   int
	Node  *data_type.RecipeTree
	Depth int
}

// Fungsi utama tanpa konkurensi, dengan batas jumlah resep
func Bfs_multiple_recipe(url string, bound int) *data_type.RecipeTree {
	idx := scrapping.MapperNameToIdx[url]
	if idx == -1 {
		fmt.Println("Error: Invalid URL")
		return nil
	}

	visited := make([]bool, 720)
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx]}
	queue := []bfsTask{{Idx: idx, Node: root, Depth: 0}}

	count := 0

	for len(queue) > 0 && count <= bound {
		task := queue[0]
		queue = queue[1:]

		idx := task.Idx
		node := task.Node
		depth := task.Depth

		if Stop(idx, visited, depth) {
			continue
		}
		visited[idx] = true

		recipes := scrapping.MapperIdxToRecipes[idx]
		for _, recipe := range recipes {
			if count >= bound {
				break
			}

			firstIdx := recipe.First
			secondIdx := recipe.Second

			firstNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[firstIdx]}
			secondNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[secondIdx]}
			pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}

			node.Children = append(node.Children, pair)
			if (firstIdx <= 4  && secondIdx <= 4) {
				count++
			}

			queue = append(queue, bfsTask{Idx: firstIdx, Node: firstNode, Depth: depth + 1})
			queue = append(queue, bfsTask{Idx: secondIdx, Node: secondNode, Depth: depth + 1})
		}
	}

	PruneNonTerminal(root)
	PruneNonTerminal(root)

	return root
}

func Stop(idx int, visited []bool, depth int) bool {
	return idx <= 4 // || visited[idx]
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

func pruneTree(node *data_type.RecipeTree) bool {
	if node == nil {
		return false
	}

	if len(node.Children) == 0 {
		return isBasicElement(node)
	}

	validChildren := []*data_type.Pair_recipe{}
	for _, pair := range node.Children {
		firstValid := pruneTree(pair.First)
		secondValid := pruneTree(pair.Second)
		if len(pair.First.Children) == 0 && len(pair.Second.Children) == 0 {
			if isBasicElement(pair.First) && isBasicElement(pair.Second) {
				validChildren = append(validChildren, pair)
			}
		} else {
			if firstValid || secondValid {
				validChildren = append(validChildren, pair)
			}
		}
	}

	node.Children = validChildren
	return len(validChildren) > 0
}


// Fungsi eksternal yang bisa kamu panggil setelah BFS selesai
func PruneNonTerminal(root *data_type.RecipeTree) {
	pruneTree(root)
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

func isBasicElement(node *data_type.RecipeTree) bool {
	idx, ok := scrapping.MapperNameToIdx[node.Name]
	return ok && idx <= 4
}
