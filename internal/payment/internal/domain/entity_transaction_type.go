package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
)

var ErrScanTransactionType = errors.New("failed to scan transaction type")

type TransactionType int

const (
	TransactionTypeUnknown TransactionType = iota
	TransactionTypeDebit
	TransactionTypeCredit
)

func ParseTransactionType(s string) TransactionType {
	sts := strings.TrimPrefix(s, "STATUS_") // for support grpc enum
	sts = strings.ToUpper(sts)

	switch sts {
	case TransactionTypeDebit.String():
		return TransactionTypeDebit
	case TransactionTypeCredit.String():
		return TransactionTypeCredit
	default:
		return TransactionTypeUnknown
	}
}

func (tt TransactionType) String() string {
	statuses := [...]string{
		"UNKNOWN",
		"DEBIT",
		"CREDIT",
	}

	if tt < TransactionTypeUnknown || int(tt) >= len(statuses) {
		return "UNKNOWN"
	}

	return statuses[tt]
}

func (tt *TransactionType) Scan(value any) error {
	switch col := value.(type) {
	case []byte:
		*tt = ParseTransactionType(string(col))
	case string:
		*tt = ParseTransactionType(col)
	default:
		return ErrScanTransactionType
	}

	return nil
}

func (tt TransactionType) Value() (driver.Value, error) {
	return tt.String(), nil
}
