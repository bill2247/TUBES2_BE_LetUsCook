package bfs_multiple_recipe

import (
	"fmt"
	"sync"
	"let_us_cook/src/scrapping"
	"let_us_cook/src/data_type"
)

// Fungsi utama untuk memulai BFS dan membentuk pohon resep secara paralel
func Bfs_multiple_recipe(url string) *data_type.RecipeTree {
	visited := make([]bool, 721)
	idx := scrapping.MapperIdxElm[url]
	if idx == -1 {
		fmt.Println("Error: Invalid URL")
		return nil
	}
	var wg sync.WaitGroup
	root := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[idx]}
	wg.Add(1)
	go bfsHelperParallel(idx, visited, 0, root, &wg)
	wg.Wait()
	return root
}

// Fungsi untuk mengecek apakah pencarian perlu dihentikan
func Stop(idx int, visited []bool, depth int) bool {
	return idx < 4 || visited[idx]
}

// Fungsi helper untuk membuat pohon resep secara paralel
func bfsHelperParallel(idx int, visited []bool, depth int, node *data_type.RecipeTree, wg *sync.WaitGroup) {
	defer wg.Done()

	if Stop(idx, visited, depth) {
		return
	}

	visited[idx] = true
	recipes := scrapping.MapperRecipe1[idx]

	for _, recipe := range recipes {
		firstIdx := recipe.First
		secondIdx := recipe.Second

		firstName := scrapping.MapperElmIdx[firstIdx]
		secondName := scrapping.MapperElmIdx[secondIdx]

		firstNode := &data_type.RecipeTree{Name: firstName}
		secondNode := &data_type.RecipeTree{Name: secondName}
		pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}
		node.Children = append(node.Children, pair)

		wg.Add(2)
		go bfsHelperParallel(firstIdx, visited, depth+1, firstNode, wg)
		go bfsHelperParallel(secondIdx, visited, depth+1, secondNode, wg)
	}
}

