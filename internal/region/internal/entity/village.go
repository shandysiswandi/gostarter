package entity

type Village struct {
	ID         string
	DistrictID string
	Name       string
}

func (v *Village) ScanColumn() []any {
	return []any{&v.ID, &v.Name}
}

func (v *Village) ToRegion() Region {
	return Region{ID: v.ID, Name: v.Name}
}
