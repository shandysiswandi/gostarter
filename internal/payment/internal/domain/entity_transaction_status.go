package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
)

var ErrScanTransactionStatus = errors.New("failed to scan transaction status")

type TransactionStatus int

const (
	TransactionStatusUnknown TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusFailed
	TransactionStatusSuccess
)

func ParseTransactionStatus(s string) TransactionStatus {
	sts := strings.TrimPrefix(s, "STATUS_") // for support grpc enum
	sts = strings.ToUpper(sts)

	switch sts {
	case TransactionStatusPending.String():
		return TransactionStatusPending
	case TransactionStatusFailed.String():
		return TransactionStatusFailed
	case TransactionStatusSuccess.String():
		return TransactionStatusSuccess
	default:
		return TransactionStatusUnknown
	}
}

func (ts TransactionStatus) String() string {
	statuses := [...]string{
		"UNKNOWN",
		"PENDING",
		"FAILED",
		"SUCCESS",
	}

	if ts < TransactionStatusUnknown || int(ts) >= len(statuses) {
		return "UNKNOWN"
	}

	return statuses[ts]
}

func (ts *TransactionStatus) Scan(value any) error {
	switch col := value.(type) {
	case []byte:
		*ts = ParseTransactionStatus(string(col))
	case string:
		*ts = ParseTransactionStatus(col)
	default:
		return ErrScanTransactionStatus
	}

	return nil
}

func (ts TransactionStatus) Value() (driver.Value, error) {
	return ts.String(), nil
}
