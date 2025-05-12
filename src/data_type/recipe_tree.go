package data_type

type RecipeTree struct {
	Name     string 		`json:"name"`
	Children []*Pair_recipe	`json:"children,omitempty"`
}