// main.go
package main

import (
	"let_us_cook/src/bfs_multiple_recipe"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "let_us_cook/src/bfs_single"
)

func main() {
	router := gin.Default()

	// Load HTML template
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")

	// Halaman utama
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/search", func(c *gin.Context) {
		startURL := c.Query("url")
		boundStr := c.Query("bound")

		bound, err := strconv.Atoi(boundStr)
		if err != nil || bound < 1 {
			bound = 5 // fallback default
		}

		tree := bfs_multiple_recipe.Bfs_multiple_recipe(startURL, bound)
		bfs_multiple_recipe.PruneNonTerminal(tree)

		// simpan ke string hasil dari DisplayTree
		treeStr := bfs_multiple_recipe.TreeToString(tree)
		c.String(http.StatusOK, treeStr) // kirim sebagai plain text
	})

	router.Run(":8080") // Jalankan server di localhost:8080
}

// func main() {
// 	url := "Duck"
// 	// max := int32(20)
// 	tree := bfs_single.FindShortestPath(url)
// 	bfs_single.DisplayTree(tree, "", true)
// }
