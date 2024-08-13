package entity

import "errors"

var (
	ErrTodoNotFound   = errors.New("todo not found")
	ErrTodoNotCreated = errors.New("todo not created")
	ErrTodoNotUpdated = errors.New("todo not updated")
	ErrTodoNotDeleted = errors.New("todo not deleted")
)

type Todo struct {
	ID          uint64
	Title       string
	Description string
	Status      TodoStatus
}

func (t *Todo) ScanColumn() []any {
	return []any{&t.ID, &t.Title, &t.Description, &t.Status}
}
