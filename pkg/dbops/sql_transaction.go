// Package dbops provides abstractions and utility functions for database operations.
// It defines interfaces for querying and executing SQL commands, and includes helper functions
// for performing common database operations, such as executing queries and scanning rows.
package dbops

import (
	"context"
	"database/sql"
	"log"
)

// contextKeySQLTx is a type used to define the context key for SQL transactions.
type contextKeySQLTx struct{}

// Transaction provides an abstraction over SQL transactions, allowing for the execution
// of a function within a transactional context.
type Transaction struct {
	db SQLTx
}

// NewTransaction creates a new Transaction instance with the provided SQLTx.
func NewTransaction(db SQLTx) *Transaction {
	return &Transaction{
		db: db,
	}
}

// Transaction executes the provided function within a transaction context.
// It starts a new transaction, attaches it to the context, and ensures that the transaction
// is either committed on success or rolled back on failure or panic.
func (t *Transaction) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, contextKeySQLTx{}, tx)

	defer func() {
		if tx == nil {
			return
		}

		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("panic when execute function and error rollback", err)
			}

			log.Println("recover from panic when execute function")

			panic(r) // is need to re-panic or not
		}

		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("error when execute function", err)
			}

			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("error when commit transaction", err)
		}
	}()

	err = fn(ctx)

	return err
}
