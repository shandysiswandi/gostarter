package domain

type District struct {
	ID     string
	CityID string
	Name   string
}

func (d *District) ScanColumn() []any {
	return []any{&d.ID, &d.Name}
}

func (d *District) ToRegion() Region {
	return Region{ID: d.ID, Name: d.Name}
}
