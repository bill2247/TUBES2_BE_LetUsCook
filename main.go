// package main

// import (
// 	// "encoding/json"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"let_us_cook/src/bfs_multiple_recipe"

// )


// // Struktur untuk request search
// type SearchRequest struct {
// 	Query     string `json:"query"`
// 	Mode      string `json:"mode"`      // 'single' atau 'multiple'
// 	Algorithm string `json:"algorithm"` // 'bfs' atau 'dfs'
// }

// func main() {
// 	r := gin.Default()

// 	// Konfigurasi CORS
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	// Route untuk search API
// 	r.POST("/api/search", func(c *gin.Context) {
// 		var req SearchRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
	
// 		// Validasi input
// 		if req.Query == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
// 			return
// 		}
// 		// ... validasi lainnya ...
	
// 		var result interface{}
		
// 		if req.Algorithm == "bfs" {
// 			if req.Mode == "multiple" {
// 				tree := bfs_multiple_recipe.Bfs_multiple_recipe(req.Query) // Gunakan req.Query sebagai URL
// 				if tree == nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
// 					return
// 				}
// 				result = tree
// 			}
// 		} else {
// 			// Handle DFS
// 		}
	
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": true,
// 			"data":    result,
// 		})
// 	})

// 	// Run server
// 	log.Println("Server is running on port 8080")
// 	r.Run(":8080")
// }

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"let_us_cook/src/bfs_multiple_recipe"
)

// Struktur untuk request search
type SearchRequest struct {
	Query     string `json:"query"`
	Mode      string `json:"mode"`      // 'single' atau 'multiple'
	Algorithm string `json:"algorithm"` // 'bfs' atau 'dfs'
}

func main() {
	r := gin.Default()

	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Simple test endpoint
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Route untuk search API
	r.POST("/api/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received search request: Query=%s, Mode=%s, Algorithm=%s", req.Query, req.Mode, req.Algorithm)

		// Validasi input
		if req.Query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
			return
		}

		if req.Mode != "single" && req.Mode != "multiple" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mode must be 'single' or 'multiple'"})
			return
		}

		if req.Algorithm != "bfs" && req.Algorithm != "dfs" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Algorithm must be 'bfs' or 'dfs'"})
			return
		}

		var result interface{}

		if req.Algorithm == "bfs" {
			if req.Mode == "multiple" {
				log.Printf("Calling BFS multiple recipe with query: %s", req.Query)
				tree := bfs_multiple_recipe.Bfs_multiple_recipe(req.Query)
				if tree == nil {
					log.Printf("BFS returned nil for query: %s", req.Query)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
					return
				}
				result = tree
			} else {
				// Handle single mode
				c.JSON(http.StatusBadRequest, gin.H{"error": "Single mode not yet implemented for BFS"})
				return
			}
		} else if req.Algorithm == "dfs" {
			// Handle DFS
			c.JSON(http.StatusBadRequest, gin.H{"error": "DFS algorithm not yet implemented"})
			return
		}

		log.Printf("Sending response with result: %+v", result)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	})

	// Run server
	log.Println("Server is running on port 8080")
	r.Run(":8080")
}