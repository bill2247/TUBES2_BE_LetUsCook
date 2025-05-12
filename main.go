package main

import (
	"encoding/json" // Tambahan
	"let_us_cook/src/bfs_multiple_recipe"
	"let_us_cook/src/bfs_shortest"
	"let_us_cook/src/dfs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"let_us_cook/src/scrapping"
)

type SearchRequest struct {
	Query     string `json:"query"`
	Mode      string `json:"mode"`
	Algorithm string `json:"algorithm"`
	CountRicipe int   `json:"countRicipe"`
}

func main() {
	scrapping.StartScraper()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://let-us-cook-new.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.OPTIONS("/api/search", func(c *gin.Context) {
    c.Status(http.StatusNoContent) 
	})

	r.POST("/api/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received search request: Query=%s, Mode=%s, Algorithm=%s, CountRicipe=%d", req.Query, req.Mode, req.Algorithm, req.CountRicipe)

		if req.Query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
			return
		}
		if req.Mode != "single" && req.Mode != "multiple" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mode must be 'single' or 'multiple'"})
			return
		}
		if req.Algorithm != "bfs" && req.Algorithm != "dfs" && req.Algorithm != "bidirectional" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Algorithm must be 'bfs', 'dfs', or 'bidirectional' "})
			return
		}

		

		var result interface{}

		if req.Algorithm == "bfs" {
			if req.Mode == "multiple" {
				log.Printf("Calling BFS multiple recipe with query: %s", req.Query)

				// hitung durasi
				startTime := time.Now()
				tree, count := bfs_multiple_recipe.Bfs_multiple_recipe(req.Query, req.CountRicipe)
				duration := time.Since(startTime)
				if tree == nil {
					log.Printf("BFS returned nil for query: %s", req.Query)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
					return
				}
				result = gin.H{
					"tree":  tree,
					"count": count, 
					"duration": duration.Seconds(),
				}

				// simpan hasil ke file JSON
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Printf("Gagal mengubah ke JSON: %v", err)
				} else {
					err := os.WriteFile("output_recipe.json", jsonBytes, 0644)
					if err != nil {
						log.Printf("Gagal menulis file JSON: %v", err)
					} else {
						log.Println("Berhasil menyimpan hasil pencarian ke output_recipe.json")
					}
				}

			} else {
				log.Printf("Calling BFS multiple recipe with query: %s", req.Query)
				// hitung durasi
				startTime := time.Now()
				tree, count := bfs_shortest.FindShortestPath(req.Query)
				duration := time.Since(startTime)

				if tree == nil {
					log.Printf("BFS returned nil for query: %s", req.Query)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
					return
				}
				result = gin.H{
					"tree":  tree,
					"count": count, 
					"duration": duration.Seconds(),
				}

				// simpan hasil ke file JSON
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Printf("Gagal mengubah ke JSON: %v", err)
				} else {
					err := os.WriteFile("output_recipe.json", jsonBytes, 0644)
					if err != nil {
						log.Printf("Gagal menulis file JSON: %v", err)
					} else {
						log.Println("Berhasil menyimpan hasil pencarian ke output_recipe.json")
					}
				}
			}
		} else if req.Algorithm == "dfs" {
			if req.Mode == "multiple" {
				log.Printf("Calling BFS multiple recipe with query: %s", req.Query)
				// hitung durasi
				startTime := time.Now()
				tree, count := dfs.DFSMultipleEntryPoint(req.Query)
				duration := time.Since(startTime)
				if tree == nil {
					log.Printf("BFS returned nil for query: %s", req.Query)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
					return
				}
				result = gin.H{
					"tree":  tree,
					"count": count, 
					"duration": duration.Seconds(),
				}
				// simpan hasil ke file JSON
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Printf("Gagal mengubah ke JSON: %v", err)
				} else {
					err := os.WriteFile("output_recipe.json", jsonBytes, 0644)
					if err != nil {
						log.Printf("Gagal menulis file JSON: %v", err)
					} else {
						log.Println("Berhasil menyimpan hasil pencarian ke output_recipe.json")
					}
				}

			} else {
				log.Printf("Calling BFS multiple recipe with query: %s", req.Query)
				// hitung durasi
				startTime := time.Now()
				tree, count := dfs.DFSSingleEntryPoint(req.Query)
				duration := time.Since(startTime)
				if tree == nil {
					log.Printf("BFS returned nil for query: %s", req.Query)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recipe"})
					return
				}
				result = gin.H{
					"tree":  tree,
					"count": count, 
					"duration": duration.Seconds(),
				}

				// simpan hasil ke file JSON
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Printf("Gagal mengubah ke JSON: %v", err)
				} else {
					err := os.WriteFile("output_recipe.json", jsonBytes, 0644)
					if err != nil {
						log.Printf("Gagal menulis file JSON: %v", err)
					} else {
						log.Println("Berhasil menyimpan hasil pencarian ke output_recipe.json")
					}
				}
			}
		} else if req.Algorithm == "bidirectional"{
			if req.Mode == "multiple"{
				// disini bidirectional
			}
		}

		log.Printf("Sending response with result: %+v", result)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": result,
		})
	})

	log.Println("Server is running on port 8080")
	r.Run(":8080")
}
