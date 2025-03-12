package domain

type Token struct {
	ID uint64 `db:"id"`
}

func (Token) Table() string {
	return "tokens"
}
