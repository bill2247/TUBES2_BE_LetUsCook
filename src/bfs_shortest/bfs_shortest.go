package bfs_shortest

import (
	"fmt"
	"let_us_cook/src/scrapping"
)

func Bfs_shortest(targetUrl string) {
	// Get index for target URL
	targetIdx := scrapping.MapperIdxElm[targetUrl]
	
	if targetIdx == -1 {
		fmt.Println("Error: Invalid target URL")
		return
	}
	
	// kasus target sudah basic element
	if targetIdx < 4 {
		fmt.Println(scrapping.MapperElmIdx[targetIdx], "is a basic element")
		return
	}
	
	parentMap, found := bfsSearch(targetIdx)
	
	if !found {
		fmt.Println("No recipe path found to", targetUrl)
		return
	}
	
	// Print the recipe tree from target to basic elements
	fmt.Println("Recipe Tree for:", scrapping.MapperElmIdx[targetIdx])
	printRecipeTree(targetIdx, parentMap, 0, make(map[int]bool))
}

// bfsSearch performs a breadth-first search to find the shortest recipe
func bfsSearch(targetIdx int) (map[int][2]int, bool) {
	// Maps to store parent relationships and distance from base elements
	parentMap := make(map[int][2]int)
	distance := make(map[int]int)
	
	// Start with basic elements (0-3) in the queue
	type QueueItem struct {
		elementIdx int
		dist       int
	}
	
	queue := []QueueItem{}
	for i := 0; i < 4; i++ {
		queue = append(queue, QueueItem{elementIdx: i, dist: 0})
		distance[i] = 0
	}
	
	// Keep track of elements we've found recipes for
	foundElements := make(map[int]bool)
	
	// BFS loop - starting from basic elements and working our way up
	for len(queue) > 0 {
		// Get next element from queue
		current := queue[0]
		queue = queue[1:]
		
		// Check if we've found the target
		if current.elementIdx == targetIdx {
			return parentMap, true
		}
		
		// Look for recipes where this element is an ingredient
		for resultIdx, recipes := range scrapping.MapperRecipe1 {
			// Skip if we already found a shorter path to this element
			if foundElements[resultIdx] {
				continue
			}
			
			for _, recipe := range recipes {
				// Check if the current element is used in this recipe
				if recipe.First == current.elementIdx || recipe.Second == current.elementIdx {
					// Find the other ingredient
					otherIdx := recipe.Second
					if otherIdx == current.elementIdx {
						otherIdx = recipe.First
					}
					
					// Check if we've processed the other ingredient
					otherDist, otherFound := distance[otherIdx]
					
					// If both ingredients are processed
					if otherFound {
						// Store the recipe
						parentMap[resultIdx] = [2]int{recipe.First, recipe.Second}
						
						// Calculate distance
						newDist := max(current.dist, otherDist) + 1
						distance[resultIdx] = newDist
						
						// Mark as found
						foundElements[resultIdx] = true
						
						// Add to queue
						queue = append(queue, QueueItem{elementIdx: resultIdx, dist: newDist})
					}
				}
			}
		}
	}
	
	// ngga ada path nya
	return parentMap, false
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// printAllCombinations prints all possible recipes until target is found
func printAllCombinations(targetIdx int) {
	// Start with basic elements
	processedElements := make(map[int]bool)
	for i := 0; i < 4; i++ {
		processedElements[i] = true
		fmt.Printf("Basic Element: %s\n", scrapping.MapperElmIdx[i])
	}
	fmt.Println()
	
	// Process recipes in levels, starting from basic elements
	level := 1
	targetFound := false
	
	// Continue until we've processed the target or can't make more progress
	for !targetFound {
		fmt.Printf("--- Level %d Combinations ---\n", level)
		levelProcessed := false
		
		// Track new elements created in this level to avoid processing them until next level
		newElementsInLevel := make(map[int]bool)
		
		// Try all possible combinations of processed elements
		for i := range processedElements {
			// Skip if element was newly created in this level
			if newElementsInLevel[i] {
				continue
			}
			
			for j := range processedElements {
				// Skip if element was newly created in this level
				if newElementsInLevel[j] {
					continue
				}
				
				// Look for recipes using this combination
				for resultIdx, recipes := range scrapping.MapperRecipe1 {
					// Skip if already processed
					if processedElements[resultIdx] {
						continue
					}
					
					// Check if this recipe uses our ingredients
					for _, recipe := range recipes {
						if (recipe.First == i && recipe.Second == j) || 
						   (recipe.First == j && recipe.Second == i) {
							// Print the combination
							fmt.Printf("%s + %s = %s\n", 
								scrapping.MapperElmIdx[i], 
								scrapping.MapperElmIdx[j], 
								scrapping.MapperElmIdx[resultIdx])
							
							// Mark as processed
							processedElements[resultIdx] = true
							newElementsInLevel[resultIdx] = true
							levelProcessed = true
							
							// Check if we found the target
							if resultIdx == targetIdx {
								targetFound = true
								fmt.Printf("  ↳ TARGET FOUND!\n")
							}
						}
					}
				}
			}
		}
		
		// If no elements were processed at this level or we found the target, we're done
		if !levelProcessed || targetFound {
			break
		}
		
		level++
		fmt.Println()
	}
	
	// Final message if target was found
	if targetFound {
		fmt.Println("\nSuccess! Target element found:", scrapping.MapperElmIdx[targetIdx])
	} else {
		fmt.Println("\nCould not find a path to target:", scrapping.MapperElmIdx[targetIdx])
	}
}

// printRecipeTree prints the recipe tree from target down to basic elements
func printRecipeTree(elementIdx int, parentMap map[int][2]int, depth int, pathTracker map[int]bool) {
	// Check for cycles
	if pathTracker[elementIdx] {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Println(indent + scrapping.MapperElmIdx[elementIdx] + " (cycle detected)")
		return
	}
	
	// Mark this element as in the current path
	pathTracker[elementIdx] = true
	
	// Get element name
	elementName := scrapping.MapperElmIdx[elementIdx]
	
	// Print element with indentation
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	// If it's a basic element, just print it
	if elementIdx < 4 {
		fmt.Println(indent + elementName + " (basic element)")
		// Unmark this element when done
		delete(pathTracker, elementIdx)
		return
	}
	
	// Print the current element
	fmt.Println(indent + elementName)
	
	// Get ingredients
	ingredients, exists := parentMap[elementIdx]
	if !exists {
		fmt.Println(indent + "  WARNING: Recipe not found")
		// Unmark this element when done
		delete(pathTracker, elementIdx)
		return
	}
	
	first := ingredients[0]
	second := ingredients[1]
	
	// Print recipes with indentation
	fmt.Println(indent + "  Recipe: " + scrapping.MapperElmIdx[first] + " + " + scrapping.MapperElmIdx[second])
	
	// Recursively print recipes for ingredients (if they're not basic elements)
	fmt.Println(indent + "  Components:")
	
	// Print the first ingredient
	printRecipeTree(first, parentMap, depth+2, pathTracker)
	
	// Print the second ingredient
	printRecipeTree(second, parentMap, depth+2, pathTracker)
	
	// Unmark this element when done with this branch
	delete(pathTracker, elementIdx)
}

// buildRecipeGraph creates a graph of all recipes needed for the target
func buildRecipeGraph(targetIdx int, parentMap map[int][2]int) map[int]bool {
	result := make(map[int]bool)
	buildRecipeGraphHelper(targetIdx, parentMap, result)
	return result
}

// buildRecipeGraphHelper is a recursive helper for buildRecipeGraph
func buildRecipeGraphHelper(elementIdx int, parentMap map[int][2]int, result map[int]bool) {
	// Mark this element as part of the recipe
	result[elementIdx] = true
	
	// If it's a basic element, we're done with this branch
	if elementIdx < 4 {
		return
	}
	
	// Get ingredients
	ingredients, exists := parentMap[elementIdx]
	if !exists {
		return
	}
	
	// Process ingredients recursively
	buildRecipeGraphHelper(ingredients[0], parentMap, result)
	buildRecipeGraphHelper(ingredients[1], parentMap, result)
}
