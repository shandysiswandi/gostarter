package http

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchResponse struct {
	Type    string   `json:"type"`
	Regions []Region `json:"regions"`
}
