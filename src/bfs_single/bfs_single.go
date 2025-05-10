package bfs_single 

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scrapping"
	"strings"
)

// fungsi utama untuk mencari jalur terpendek
func FindShortestPath(targetURL string) *data_type.RecipeTree {
	targetIdx, ok := scrapping.MapperIdxElm[targetURL]
	if !ok || targetIdx == -1 {
		fmt.Println("Error: Invalid target URL")
		return nil
	}

	// queue untuk BFS - menyimpan index resep dan tree node
	type queueItem struct {
		idx      int
		node     *data_type.RecipeTree
		distance int
	}

	// map untuk menyimpan node yang sudah dikunjungi dan jaraknya
	visited := make(map[int]int) 
	
	root := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[targetIdx]}
	
	// inisialisasi queue dengan target resep
	queue := []queueItem{{idx: targetIdx, node: root, distance: 0}}
	visited[targetIdx] = 0
	
	// map untuk menyimpan resep terbaik untuk setiap elemen
	bestRecipes := make(map[int]*data_type.Recipe)
	
	// BFS untuk mencari jalur terpendek
	for len(queue) > 0 {
		// ambil item pertama dari queue
		current := queue[0]
		queue = queue[1:]
		
		currentIdx := current.idx
		currentNode := current.node
		currentDistance := current.distance
		
		// jika sudah mencapai elemen dasar, stop
		if currentIdx <= 4 {
			continue
		}
		
		// semua resep yang bisa membuat elemen saat ini
		recipes := scrapping.MapperRecipe1[currentIdx]
		
		if len(recipes) == 0 {
			continue
		}
		
		// temukan resep terbaik (yang menggunakan elemen paling dasar)
		var bestRecipe *data_type.Recipe
		bestScore := -1
		
		for _, recipe := range recipes {
			firstIdx := recipe.First
			secondIdx := recipe.Second
			
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
		
		// buat node untuk bahan pertama dan kedua
		firstNode := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[firstIdx]}
		secondNode := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[secondIdx]}
		
		// tambahkan ke dalam tree
		pair := &data_type.Pair_recipe{First: firstNode, Second: secondNode}
		currentNode.Children = append(currentNode.Children, pair)
		
		// tambahkan ke queue untuk penelusuran lebih lanjut jika belum dikunjungi atau jalur lebih pendek
		newDistance := currentDistance + 1
		
		if dist, found := visited[firstIdx]; !found || newDistance < dist {
			visited[firstIdx] = newDistance
			queue = append(queue, queueItem{idx: firstIdx, node: firstNode, distance: newDistance})
		}
		
		if dist, found := visited[secondIdx]; !found || newDistance < dist {
			visited[secondIdx] = newDistance
			queue = append(queue, queueItem{idx: secondIdx, node: secondNode, distance: newDistance})
		}
	}
	
	// rekonstruksi jalur optimal
	optimalPath := ConstructOptimalPath(targetIdx, bestRecipes)
	
	return optimalPath
}

// rekonstruksi jalur optimal dari resep terbaik
func ConstructOptimalPath(targetIdx int, bestRecipes map[int]*data_type.Recipe) *data_type.RecipeTree {
	root := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[targetIdx]}
	
	// DFS untuk membangun pohon resep
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
	
	// Dapatkan resep terbaik untuk elemen ini
	recipe, found := bestRecipes[idx]
	if !found {
		return
	}
	
	// buat node untuk kedua bahan
	firstIdx := recipe.First
	secondIdx := recipe.Second
	
	firstNode := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[firstIdx]}
	secondNode := &data_type.RecipeTree{Name: scrapping.MapperElmIdx[secondIdx]}
	
	// tambahkan ke dalam tree
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

func isBasicElement(node *data_type.RecipeTree) bool {
	idx, ok := scrapping.MapperIdxElm[node.Name]
	return ok && idx <= 4
}

// cut cabang yang tidak mengarah ke elemen dasar
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
		
		// jika keduanya elemen dasar atau keduanya valid
		if (isBasicElement(pair.First) && isBasicElement(pair.Second)) || 
		   (firstValid && secondValid) {
			validChildren = append(validChildren, pair)
		}
	}
	
	node.Children = validChildren
	return len(validChildren) > 0
}

// fungsi eksternal untuk memangkas tree
func PruneNonTerminal(root *data_type.RecipeTree) {
	pruneTree(root)
}

// menampilkan tree ke konsol
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

// konversi tree ke string
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

// mengekstrak langkah-langkah resep dan mengembalikan dalam bentuk string
func GetRecipeSteps(root *data_type.RecipeTree) string {
	if root == nil {
		return "Resep tidak ditemukan"
	}
	
	var steps []string
	extractSteps(root, &steps)
	
	result := "Langkah-langkah membuat " + root.Name + ":\n"
	for i, step := range steps {
		result += fmt.Sprintf("%d. %s\n", i+1, step)
	}
	
	return result
}

// langkah-langkah resep dari tree
func extractSteps(node *data_type.RecipeTree, steps *[]string) {
	if node == nil || len(node.Children) == 0 {
		return
	}
	
	for _, pair := range node.Children {
		// tambahkan langkah resep
		step := fmt.Sprintf("%s = %s + %s", node.Name, pair.First.Name, pair.Second.Name)
		*steps = append(*steps, step)
		
		// rekursif untuk bahan-bahan
		extractSteps(pair.First, steps)
		extractSteps(pair.Second, steps)
	}
}