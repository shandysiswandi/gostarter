package domain

type User struct {
	ID       uint64
	Name     string
	Email    string
	Password string
}

func (u *User) ScanColumn() []any {
	return []any{&u.ID, &u.Name, &u.Email, &u.Password}
}
