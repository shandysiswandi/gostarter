package domain

import (
	"errors"

	"github.com/shandysiswandi/gostarter/pkg/enum"
)

var (
	ErrTodoNotFound   = errors.New("todo not found")
	ErrTodoNotCreated = errors.New("todo not created")
	ErrTodoNotUpdated = errors.New("todo not updated")
	ErrTodoNotDeleted = errors.New("todo not deleted")
)

type Todo struct {
	ID          uint64
	UserID      uint64
	Title       string
	Description string
	Status      enum.Enum[TodoStatus]
}

func (t *Todo) ScanColumn() []any {
	return []any{&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status}
}
