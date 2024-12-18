package domain

import (
	"errors"
	"time"

	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shopspring/decimal"
)

var ErrTransactionNoRowsAffected = errors.New("transaction not created or update")

type TransactionStatus int

const (
	TransactionStatusUnknown TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusFailed
	TransactionStatusSuccess
)

func (ts TransactionStatus) Values() map[enum.Enumerate]string {
	return map[enum.Enumerate]string{
		TransactionStatusUnknown: "UNKNOWN",
		TransactionStatusPending: "PENDING",
		TransactionStatusFailed:  "FAILED",
		TransactionStatusSuccess: "SUCCESS",
	}
}

type TransactionType int

const (
	TransactionTypeUnknown TransactionType = iota
	TransactionTypeDebit
	TransactionTypeCredit
)

func (tt TransactionType) Values() map[enum.Enumerate]string {
	return map[enum.Enumerate]string{
		TransactionStatusUnknown: "UNKNOWN",
		TransactionTypeDebit:     "DEBIT",
		TransactionTypeCredit:    "CREDIT",
	}
}

type Transaction struct {
	ID       uint64
	UserID   uint64
	Amount   decimal.Decimal
	Type     TransactionType
	Status   TransactionStatus
	Remark   string
	CreateAt time.Time
}

func (t *Transaction) ScanColumn() []any {
	return []any{
		&t.ID,
		&t.UserID,
		&t.Amount,
		&t.Type,
		&t.Status,
		&t.Remark,
		&t.CreateAt,
	}
}
