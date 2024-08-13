package outbound

import (
	"context"
	"database/sql"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/pkg/persistence"
)

type MysqlTodo struct {
	db *sql.DB
}

func NewMysqlTodo(db *sql.DB) *MysqlTodo {
	return &MysqlTodo{
		db: db,
	}
}

func (mt *MysqlTodo) Create(ctx context.Context, todo entity.Todo) error {
	query := `INSERT INTO todos(id,title,description,status) VALUES(?,?,?,?)`

	ir, err := mt.db.ExecContext(ctx, query,
		todo.ID, todo.Title, todo.Description, todo.Status.String())
	if err != nil {
		return err
	}

	aff, err := ir.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return entity.ErrTodoNotCreated
	}

	return nil
}

func (mt *MysqlTodo) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM todos WHERE id=?`

	ir, err := mt.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	aff, err := ir.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return entity.ErrTodoNotDeleted
	}

	return nil
}

func (mt *MysqlTodo) GetByID(ctx context.Context, id uint64) (*entity.Todo, error) {
	query := func() (string, []any, error) {
		query := "SELECT id, title, description,status FROM todos WHERE id = ?"

		return query, []any{id}, nil
	}

	todo, err := persistence.SQLGet[entity.Todo](ctx, mt.db, query)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (mt *MysqlTodo) GetWithFilter(ctx context.Context, _ map[string]string) ([]entity.Todo, error) {
	query := func() (string, []any, error) {
		return "SELECT id, title, description, status FROM todos", nil, nil
	}

	todos, err := persistence.SQLGets[entity.Todo](ctx, mt.db, query)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (mt *MysqlTodo) UpdateStatus(ctx context.Context, id uint64, sts entity.TodoStatus) error {
	query := `UPDATE todos SET status=? WHERE id=?`

	ir, err := mt.db.ExecContext(ctx, query, sts.String(), id)
	if err != nil {
		return err
	}

	aff, err := ir.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return entity.ErrTodoNotUpdated
	}

	return nil
}

func (mt *MysqlTodo) Update(ctx context.Context, todo entity.Todo) error {
	query := `UPDATE todos SET title=?,description=?,status=? WHERE id=?`

	ir, err := mt.db.ExecContext(ctx, query,
		todo.Title, todo.Description, todo.Status.String(), todo.ID)
	if err != nil {
		return err
	}

	aff, err := ir.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return entity.ErrTodoNotUpdated
	}

	return nil
}
