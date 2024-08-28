package domain

import (
	"time"
)

type Short struct {
	Key     string
	URL     string
	Slug    string
	Expired time.Time
}
