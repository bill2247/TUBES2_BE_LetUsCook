// main.go
package main

// import (
// 	"bufio"
// 	"fmt"
// 	"let_us_cook/src/dfs"
// 	"let_us_cook/src/scrapping"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"
// )

// func main() {
// 	fmt.Println("[Kusanagi Nene's Brewery]")
// 	reader := bufio.NewReader(os.Stdin)
// 	input, _ := reader.ReadString('\n')
// 	root := dfs.CreateRecipeTreeFromName(strings.TrimSpace(input))
// 	dfs.GlobalCounter.SetCounter(scrapping.MapperIdxElm[root.Name], 1000)

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	start := time.Now()
// 	go func() {
// 		dfs.DFSMultiple(root, &wg, 0)
// 	}()
// 	wg.Wait()
// 	elapsed := time.Since(start)

// 	fmt.Println("------------------------")
// 	dfs.PrintTree(root, 0)
// 	fmt.Printf("Execution time: %s\n", elapsed)
// 	fmt.Printf("Node count    : %d\n", dfs.NodeCount(root))
// 	fmt.Printf("Found         : %d\n", dfs.GlobalCounter.GetCount())

// 	dfs.GlobalCounter.SetCounter(scrapping.MapperIdxElm[root.Name], 1000)

// 	start2 := time.Now()
// 	dfs.DFSMultipleSerial(root)
// 	elapsed2 := time.Since(start2)

// 	fmt.Println("------------------------")
// 	dfs.PrintTree(root, 0)
// 	fmt.Printf("Execution time: %s\n", elapsed2)
// 	fmt.Printf("Node count    : %d\n", dfs.NodeCount(root))
// 	fmt.Printf("Found         : %d\n", dfs.GlobalCounter.GetCount())
// }
