package http

import "github.com/shandysiswandi/gostarter/internal/region/internal/entity"

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FromListEntityRegion(list []entity.Region) []Region {
	result := make([]Region, 0)
	for _, r := range list {
		result = append(result, Region{ID: r.ID, Name: r.Name})
	}

	return result
}

// SearchResponse is ..
type SearchResponse struct {
	Regions []Region `json:"regions"`
}
