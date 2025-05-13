package scrapping

import (
	"errors"
	"fmt"
	"let_us_cook/src/data_type"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	FinalData          = make(map[string][]data_type.AlchemyEntry)
	MapperNameToIdx    = make(map[string]int)
	MapperIdxToName    = make(map[int]string)
	MapperIdxToTier    = make(map[int]int)
	MapperIdxToRecipes = make(map[int][]data_type.Recipe)
	MapperPairToIdxs  = make(map[data_type.Recipe][]int)
)

// normalizeText membersihkan spasi berlebih dari teks
func normalizeText(input string) string {
	pattern := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(pattern.ReplaceAllString(input, " "))
}

// extractCombinations mengurai cell tabel menjadi daftar resep
func extractCombinations(td *goquery.Selection) []string {
	var results []string

	td.Find("li").Each(func(_ int, item *goquery.Selection) {
		content := normalizeText(item.Text())
		if content != "" {
			results = append(results, content)
		}
	})

	// fallback jika tidak ada <li>
	if len(results) == 0 {
		fallback := normalizeText(td.Text())
		if fallback != "" {
			results = append(results, fallback)
		}
	}

	return results
}

// StartScraper menjalankan proses scraping utama
func StartScraper() error {
	targetURL := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	response, err := http.Get(targetURL)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status code: %d %s", response.StatusCode, response.Status)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return fmt.Errorf("HTML parsing error: %w", err)
	}

	idxCounter := 0

	doc.Find("h3").Each(func(_ int, header *goquery.Selection) {
		span := header.Find("span.mw-headline")
		if span.Length() == 0 {
			return
		}

		category := normalizeText(span.Text())
		FinalData[category] = []data_type.AlchemyEntry{}

		// Tentukan tier
		var tier int
		switch strings.ToLower(category) {
		case "special element":
			tier = -1
		case "starting elements":
			tier = 0
		default:
			_, err := fmt.Sscanf(strings.ToLower(category), "tier %d elements", &tier)
			if err != nil {
				tier = -999 // fallback tier
			}
		}

		// Cari tabel berikutnya setelah h3
		tableNode := header.Next()
		for tableNode != nil && goquery.NodeName(tableNode) != "table" {
			tableNode = tableNode.Next()
		}
		if tableNode == nil {
			return
		}

		tableNode.Find("tr").Each(func(rowIdx int, row *goquery.Selection) {
			if rowIdx == 0 {
				return // skip header
			}

			cells := row.Find("td")
			if cells.Length() >= 2 {
				elementName := normalizeText(cells.Eq(0).Text())
				recipeList := extractCombinations(cells.Eq(1))

				FinalData[category] = append(FinalData[category], data_type.AlchemyEntry{
					Name:     elementName,
					Combines: recipeList,
				})

				// Isi mapper jika belum ada
				if _, exists := MapperNameToIdx[elementName]; !exists {
					MapperNameToIdx[elementName] = idxCounter
					MapperIdxToName[idxCounter] = elementName
					MapperIdxToTier[idxCounter] = tier
					idxCounter++
				}
			}
		})
	})

	// Bangun MapperIdxToRecipes setelah FinalData selesai
	for _, entries := range FinalData {
		for _, entry := range entries {
			resultIdx, ok := MapperNameToIdx[entry.Name]
			if !ok {
				continue
			}
			for _, recipe := range entry.Combines {
				parts := strings.Split(recipe, "+")
				if len(parts) != 2 {
					continue
				}
				first := normalizeText(parts[0])
				second := normalizeText(parts[1])
				if first == "Time" || second == "Time" {
					continue
				}
				firstIdx, ok1 := MapperNameToIdx[first]
				secondIdx, ok2 := MapperNameToIdx[second]
				pair1 := data_type.Recipe{First: firstIdx, Second: secondIdx}
				pair2 := data_type.Recipe{First: secondIdx, Second: firstIdx}
				if ok1 && ok2 {
					MapperIdxToRecipes[resultIdx] = append(MapperIdxToRecipes[resultIdx], pair1)
					MapperPairToIdxs[pair1] = append(MapperPairToIdxs[pair1], resultIdx)
					MapperPairToIdxs[pair2] = append(MapperPairToIdxs[pair2], resultIdx)
					
				}
			}
		}
	}
	return errors.New("scraping completed")
}
