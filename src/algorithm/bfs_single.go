package algorithm 

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scraping"
)

// fungsi utama untuk mencari jalur terpendek
func FindShortestPath(targetURL string) (*data_type.RecipeTree, int) {
	targetIdx, ok := scrapping.MapperNameToIdx[targetURL]
	if !ok || targetIdx == -1 {
		fmt.Println("Error: Invalid target URL")
		return nil, 0
	}

	// queue untuk BFS - menyimpan index resep dan tree node
	type queueItem struct {
		idx      int
		node     *data_type.RecipeTree
		distance int
		path     map[int]bool // track elemen dalam path saat ini
	}
	tier := scrapping.MapperIdxToTier[targetIdx]
	if tier == -1 || tier == 0 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[targetIdx], Children: nil}, 1
	}

	// map untuk menyimpan node yang sudah dikunjungi dan jaraknya
	visited := make(map[int]int) 
	
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[targetIdx]}
	
	// inisialisasi queue dengan target resep
	initialPath := make(map[int]bool)
	initialPath[targetIdx] = true // tandai target sebagai bagian dari path
	
	queue := []queueItem{{idx: targetIdx, node: root, distance: 0, path: initialPath}}
	visited[targetIdx] = 0
	
	// map untuk menyimpan resep terbaik untuk setiap elemen
	bestRecipes := make(map[int]*data_type.Recipe)
	
	// BFS untuk mencari jalur terpendek
	visitedCount := 1
	for len(queue) > 0 {
		// ambil item pertama dari queue
		current := queue[0]
		queue = queue[1:]
		
		currentIdx := current.idx
		currentNode := current.node
		currentDistance := current.distance
		currentPath := current.path
		
		// jika sudah mencapai elemen dasar, stop
		if currentIdx <= 4 {
			continue
		}
		
		// semua resep yang bisa membuat elemen saat ini
		recipes := scrapping.MapperIdxToRecipes[currentIdx]
		
		if len(recipes) == 0 {
			continue
		}
		
		// temukan resep terbaik (yang menggunakan elemen dari tier yang lebih rendah)
		var bestRecipe *data_type.Recipe
		bestScore := -1
		
		currentTier := scrapping.MapperIdxToTier[currentIdx]
		
		for _, recipe := range recipes {
			firstIdx := recipe.First
			secondIdx := recipe.Second
			
			firstTier := scrapping.MapperIdxToTier[firstIdx]
			secondTier := scrapping.MapperIdxToTier[secondIdx]
			
			// pastikan kedua bahan berasal dari tier yang lebih rendah
			if firstTier >= currentTier || secondTier >= currentTier {
				continue
			}
			
			// mencegah siklus dan mencegah menggunakan elemen target sebagai bahan
			if currentPath[firstIdx] || currentPath[secondIdx] {
				continue
			}
			
			// hitung skor resep (prioritaskan elemen dasar)
			score := 0
			if firstIdx <= 4 {
				score++
			}
			if secondIdx <= 4 {
				score++
			}
			
			// pilih resep dengan skor tertinggi atau yang belum pernah dikunjungi
			if score > bestScore || bestRecipe == nil {
				bestRecipe = &recipe
				bestScore = score
			}
		}
		
		if bestRecipe == nil {
			continue
		}
		
		bestRecipes[currentIdx] = bestRecipe
		
		firstIdx := bestRecipe.First
		secondIdx := bestRecipe.Second
		
		visitedCount += 2
		// buat node untuk bahan pertama dan kedua
		firstNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[firstIdx]}
		secondNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[secondIdx]}
		
		// tambahkan ke dalam tree
		pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}
		currentNode.Children = append(currentNode.Children, pair)
		
		newDistance := currentDistance + 1
		
		// buat copy dari path saat ini untuk digunakan di jalur baru
		newPathFirst := make(map[int]bool)
		for k, v := range currentPath {
			newPathFirst[k] = v
		}
		newPathFirst[firstIdx] = true
		
		newPathSecond := make(map[int]bool)
		for k, v := range currentPath {
			newPathSecond[k] = v
		}
		newPathSecond[secondIdx] = true
		
		if dist, found := visited[firstIdx]; !found || newDistance < dist {
			visited[firstIdx] = newDistance
			queue = append(queue, queueItem{idx: firstIdx, node: firstNode, distance: newDistance, path: newPathFirst})
		}
		
		if dist, found := visited[secondIdx]; !found || newDistance < dist {
			visited[secondIdx] = newDistance
			queue = append(queue, queueItem{idx: secondIdx, node: secondNode, distance: newDistance, path: newPathSecond})
		}
	}
	
	// rekonstruksi jalur optimal
	optimalPath := ConstructOptimalPath(targetIdx, bestRecipes)
	return optimalPath, visitedCount
}

func ConstructOptimalPath(targetIdx int, bestRecipes map[int]*data_type.Recipe) *data_type.RecipeTree {
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[targetIdx]}
	
	BuildRecipeTree(root, targetIdx, bestRecipes, make(map[int]bool))
	
	return root
}

// membangun pohon resep rekursif
func BuildRecipeTree(node *data_type.RecipeTree, idx int, bestRecipes map[int]*data_type.Recipe, visited map[int]bool) {
	// hentikan jika sudah mencapai elemen dasar atau sudah dikunjungi (menghindari loop)
	if idx <= 4 || visited[idx] {
		return
	}
	
	visited[idx] = true
	
	recipe, found := bestRecipes[idx]
	if !found {
		return
	}
	
	// buat node untuk kedua bahan
	firstIdx := recipe.First
	secondIdx := recipe.Second
	
	firstNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[firstIdx]}
	secondNode := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[secondIdx]}
	
	pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}
	node.Children = append(node.Children, pair)
	
	// rekursif untuk bahan-bahan jika bukan elemen dasar
	if firstIdx > 4 {
		BuildRecipeTree(firstNode, firstIdx, bestRecipes, visited)
	}
	
	if secondIdx > 4 {
		BuildRecipeTree(secondNode, secondIdx, bestRecipes, visited)
	}
}