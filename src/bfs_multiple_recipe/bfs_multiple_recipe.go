package bfs_multiple_recipe

import (
	"fmt"
	"src/scrapping"
)

func Bfs_multiple_recipe(url string) {
	visited := make([]bool, 721)

	idx := scrapping.MapperIdxElm[url]
	if idx == -1 {
		fmt.Println("Error: Invalid URL")
		return
	}
	Bfs_helper(idx, visited, 0)
}

func Stop(idx int, visited []bool) bool {
	return idx < 4 || visited[idx]
}

func Bfs_helper(idx int, visited []bool, depth int) {
	// Buat indentasi berdasarkan depth
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "--"
	}

	this := scrapping.MapperElmIdx[idx]
	fmt.Println(indent + this)

	if Stop(idx, visited) {
		return
	}

	visited[idx] = true
	recipes := scrapping.MapperRecipe1[idx]

	for _, recipe := range recipes {
		Bfs_helper(recipe.First, visited, depth+1)
		Bfs_helper(recipe.Second, visited, depth+1)
	}
}
