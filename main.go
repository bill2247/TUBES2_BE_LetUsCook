// main.go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"let_us_cook/src/bfs_multiple_recipe" 
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

	// Endpoint untuk memulai BFS
	router.GET("/search", func(c *gin.Context) {
		startURL := c.Query("url")
		tree := bfs_multiple_recipe.Bfs_multiple_recipe(startURL)
		c.JSON(http.StatusOK, tree) // kirim sebagai JSON
	})

	router.Run(":8080") // Jalankan server di localhost:8080
}
