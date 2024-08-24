package entity

type Province struct {
	ID   string
	Name string
}

func (p *Province) ScanColumn() []any {
	return []any{&p.ID, &p.Name}
}

func (p *Province) ToRegion() Region {
	return Region{ID: p.ID, Name: p.Name}
}
