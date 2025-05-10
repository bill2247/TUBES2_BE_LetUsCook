package scrapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"let_us_cook/src/data_type"

	"github.com/PuerkitoBio/goquery"
)

var (
	FinalData  = make(map[string][]AlchemyEntry)
	MapperNameToIdx    = make(map[string]int)
	MapperIdxToName    = make(map[int]string)
	MapperIdxToTier    = make(map[int]int)
	MapperIdxToRecipes = make(map[int][]data_type.Recipe)
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

// AlchemyEntry merepresentasikan elemen dan daftar resepnya
type AlchemyEntry struct {
	Name     string   `json:"element"`
	Combines []string `json:"recipes"`
}

// FetchAllData memuat semua data dari file JSON ke slice
func FetchAllData() ([]AlchemyEntry, error) {
	source, err := os.Open("little_alchemy_elements.json")
	if err != nil {
		return nil, errors.New("failed to open saved data file")
	}
	defer source.Close()

	container := make(map[string][]AlchemyEntry)
	if err := json.NewDecoder(source).Decode(&container); err != nil {
		return nil, errors.New("JSON decoding failed")
	}

	var flatList []AlchemyEntry
	for _, group := range container {
		flatList = append(flatList, group...)
	}
	return flatList, nil
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
		FinalData[category] = []AlchemyEntry{}

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

				FinalData[category] = append(FinalData[category], AlchemyEntry{
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
				firstIdx, ok1 := MapperNameToIdx[first]
				secondIdx, ok2 := MapperNameToIdx[second]
				pair :=  data_type.Recipe{First: firstIdx, Second: secondIdx}
				if ok1 && ok2 {
					MapperIdxToRecipes[resultIdx] = append(MapperIdxToRecipes[resultIdx], pair)
				}
			}
		}
	}

	FinalDataSaveToFile()
	MapperNameToIdxSaveToFile()
	MapperIdxToNameSaveToFile()
	MapperIdxToTierSaveToFile()
	MapperIdxToRecipesSaveToFile()

	return errors.New("scraping completed")
}

func FinalDataSaveToFile() {
	// Simpan FinalData ke file JSON
	output, err := os.Create("scraper/JSON/little_alchemy_elements.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer output.Close()

	writer := json.NewEncoder(output)
	writer.SetIndent("", "  ")
	if err := writer.Encode(FinalData); err != nil {
		fmt.Println("Error writing JSON content:", err)
		return
	}
	fmt.Println("Final data saved to little_alchemy_elements.json")
}

func MapperNameToIdxSaveToFile(){
	// Simpan MapperNameToIdx ke file JSON
	output, err := os.Create("scraper/JSON/MapperNameToIdx.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer output.Close()

	writer := json.NewEncoder(output)
	writer.SetIndent("", "  ")
	if err := writer.Encode(MapperNameToIdx); err != nil {
		fmt.Println("Error writing JSON content:", err)
		return
	}
	fmt.Println("Mapper saved to MapperNameToIdx.json")
}

func MapperIdxToNameSaveToFile(){
	// Simpan MapperIdxToName ke file JSON
	output, err := os.Create("scraper/JSON/MapperIdxToName.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer output.Close()

	writer := json.NewEncoder(output)
	writer.SetIndent("", "  ")
	if err := writer.Encode(MapperIdxToName); err != nil {
		fmt.Println("Error writing JSON content:", err)
		return
	}
	fmt.Println("Mapper saved to MapperIdxToName.json")
}

func MapperIdxToTierSaveToFile(){
	// Simpan MapperIdxToTier ke file JSON
	output, err := os.Create("scraper/JSON/MapperIdxToTier.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer output.Close()

	writer := json.NewEncoder(output)
	writer.SetIndent("", "  ")
	if err := writer.Encode(MapperIdxToTier); err != nil {
		fmt.Println("Error writing JSON content:", err)
		return
	}
	fmt.Println("Mapper saved to MapperIdxToTier.json")
}

func MapperIdxToRecipesSaveToFile(){
	// Simpan MapperIdxToRecipes ke file JSON
	output, err := os.Create("scraper/JSON/MapperIdxToRecipes.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer output.Close()

	writer := json.NewEncoder(output)
	writer.SetIndent("", "  ")
	if err := writer.Encode(MapperIdxToRecipes); err != nil {
		fmt.Println("Error writing JSON content:", err)
		return
	}
	fmt.Println("Mapper saved to MapperIdxToRecipes.json")
}