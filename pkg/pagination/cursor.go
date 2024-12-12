package pagination

import (
	"encoding/base64"
	"strconv"
)

const DefaultLimit = 10
const MaxLimit = 100

func ParseCursorBased(cursor, limit string) (uint64, int) {
	lmt, err := strconv.Atoi(limit)
	if err != nil {
		lmt = DefaultLimit
	}

	if lmt <= 0 {
		lmt = DefaultLimit
	}

	if lmt > MaxLimit {
		lmt = MaxLimit
	}

	if cursor == "" {
		return 0, lmt
	}

	cursorBytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, lmt
	}

	csr, err := strconv.ParseUint(string(cursorBytes), 10, 64)
	if err != nil {
		return 0, lmt
	}

	return csr, lmt
}
