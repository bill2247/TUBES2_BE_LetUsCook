package data_type

// AlchemyEntry merepresentasikan elemen dan daftar resepnya
type AlchemyEntry struct {
	Name     string   `json:"element"`
	Combines []string `json:"recipes"`
}