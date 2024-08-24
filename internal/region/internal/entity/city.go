package entity

type City struct {
	ID         string
	ProvinceID string
	Name       string
}

func (c *City) ScanColumn() []any {
	return []any{&c.ID, &c.ProvinceID, &c.Name}
}

func (c *City) ToRegion() Region {
	return Region{ID: c.ID, Name: c.Name}
}
