package dfs

import (
	"fmt"
	"let_us_cook/src/data_type"
	"let_us_cook/src/scrapping"
	"sync"
)

type Counter struct {
	rootId int
	mu     sync.Mutex
	count  int
	limit  int
}

var GlobalCounter = Counter{rootId: 0, count: 0, limit: 0}

// var MAX_CONCURRENCY	= 16 // semoga tidak black screen
// var semaphore = make(chan struct{}, MAX_CONCURRENCY)

func (c *Counter) SetCounter(id int, limit int) {
	c.rootId = id
	c.count = 0
	c.limit = limit
}

func (c *Counter) GetCount() int {
	return c.count
}

func (c *Counter) TryAdd(n int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.count >= c.limit {
		return false
	}
	c.count += n
	return c.count < c.limit
}

func (c *Counter) IsLimitReached() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count >= c.limit
}

func DFSMultiple(t *data_type.RecipeTree, wg *sync.WaitGroup) int {
	if wg != nil {
		defer wg.Done()
	}

	currentId := scrapping.MapperIdxElm[t.Name]

	// basis
	if currentId == 0 || currentId == 1 || currentId == 2 || currentId == 3 {
		return 1
	}

	// rekurens
	childrenList := scrapping.MapperRecipe1[currentId]

	// multithreading
	totalWays := 0
	for i := 0; i < len(childrenList); i++ {
		if i > 0 && GlobalCounter.IsLimitReached() {
			return 1
		}
		if totalWays >= GlobalCounter.limit {
			return totalWays
		}

		idFirst := childrenList[i].First
		idSecond := childrenList[i].Second
		if getTier(idFirst) >= getTier(currentId) || getTier(idSecond) >= getTier(currentId) {
			continue
		}

		firstRecipe := CreateRecipeTreeFromId(idFirst)
		secondRecipe := CreateRecipeTreeFromId(idSecond)
		currentPair := &data_type.Pair_recipe{First: firstRecipe, Second: secondRecipe}
		t.Children = append(t.Children, currentPair)

		channel1 := make(chan int, 1)
		channel2 := make(chan int, 1)

		var childWg sync.WaitGroup
		childWg.Add(2)

		go func() {
			defer childWg.Done()
			channel1 <- DFSMultiple(firstRecipe, nil)
		}()
		go func() {
			defer childWg.Done()
			channel2 <- DFSMultiple(secondRecipe, nil)
		}()

		childWg.Wait()

		countFirst := <-channel1
		countSecond := <-channel2
		totalWays += countFirst * countSecond

		if currentId == GlobalCounter.rootId {
			GlobalCounter.TryAdd(countFirst * countSecond)
		}
	}
	return totalWays
}

func DFSMultipleEntryPoint(url string) *data_type.RecipeTree {
	idx := scrapping.MapperNameToIdx[url]
	if idx == -1 {
		fmt.Println("Error: Invalid URL")
		return nil
	}
	tier := scrapping.MapperIdxToTier[idx]
	if tier == -1 {
		return &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx], Children: nil}
	}
	root := &data_type.RecipeTree{Name: scrapping.MapperIdxToName[idx]}
	GlobalCounter.SetCounter(scrapping.MapperIdxElm[root.Name], 100)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		DFSMultiple(root, &wg)
	}()
	wg.Wait()
	return root
}

// for debugging purposes
func DFSMultipleSerial(t *data_type.RecipeTree) int {
	if GlobalCounter.IsLimitReached() {
		return 0
	}

	currentId := scrapping.MapperIdxElm[t.Name]

	// basis
	if currentId == 0 || currentId == 1 || currentId == 2 || currentId == 3 {
		return 1
	}

	// rekurens
	childrenList := scrapping.MapperRecipe1[currentId]
	totalWays := 0
	for i := 0; i < len(childrenList); i++ {
		if i > 0 && GlobalCounter.IsLimitReached() {
			return 1
		}
		if totalWays >= GlobalCounter.limit {
			return totalWays
		}
		idFirst := childrenList[i].First
		idSecond := childrenList[i].Second

		// skip yang tidak memenuhi kriteria
		if getTier(idFirst) >= getTier(currentId) || getTier(idSecond) >= getTier(currentId) {
			continue
		}

		firstRecipe := CreateRecipeTreeFromId(idFirst)
		secondRecipe := CreateRecipeTreeFromId(idSecond)

		currentPair := &data_type.Pair_recipe{First: firstRecipe, Second: secondRecipe}
		t.Children = append(t.Children, currentPair)

		countFirst := DFSMultipleSerial(firstRecipe)
		countSecond := DFSMultipleSerial(secondRecipe)

		totalWays += countFirst * countSecond

		if currentId == GlobalCounter.rootId {
			GlobalCounter.TryAdd(countFirst * countSecond)
		}
	}
	return totalWays
}
