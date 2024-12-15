package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
)

var ErrScanBillType = errors.New("failed to scan bill type")

type BillType int

const (
	BillTypeUnknown BillType = iota
	BillTypePulsa
	BillTypeListrik
	BillTypeInternet
)

func ParseBillType(s string) BillType {
	sts := strings.TrimPrefix(s, "STATUS_") // for support grpc enum
	sts = strings.ToUpper(sts)

	switch sts {
	case BillTypePulsa.String():
		return BillTypePulsa
	case BillTypeListrik.String():
		return BillTypeListrik
	case BillTypeInternet.String():
		return BillTypeInternet
	default:
		return BillTypeUnknown
	}
}

func (bt BillType) String() string {
	statuses := [...]string{
		"UNKNOWN",
		"PULSA",
		"LISTRIK",
		"INTERNET",
	}

	if bt < BillTypeUnknown || int(bt) >= len(statuses) {
		return "UNKNOWN"
	}

	return statuses[bt]
}

func (bt *BillType) Scan(value any) error {
	switch col := value.(type) {
	case []byte:
		*bt = ParseBillType(string(col))
	case string:
		*bt = ParseBillType(col)
	default:
		return ErrScanBillType
	}

	return nil
}

func (bt BillType) Value() (driver.Value, error) {
	return bt.String(), nil
}
