package domain

type User struct {
	ID       uint64
	Email    string
	Password string
}

func (u *User) ScanColumn() []any {
	return []any{&u.ID, &u.Email, &u.Password}
}
