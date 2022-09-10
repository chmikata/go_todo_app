package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/chmikata/go_todo_app/entity"
	"github.com/lib/pq"
)

func (r *Repository) RegisterUser(
	ctx context.Context, db Execer, u *entity.User,
) error {
	sql := `insert into todoapp.users(
		name, password, role, created, modified
	) values ($1, $2, $3, $4, $5) returning id;`
	err := db.QueryRowxContext(
		ctx, sql, u.Name, u.Password,
		u.Role, r.Clocker.Now(), r.Clocker.Now(),
	).Scan(&u.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pq.ErrorCode(ErrCodePqSQLDuplicateEntry) {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}
	return nil
}

func (r *Repository) GetUser(
	ctx context.Context, db Queryer, name string,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `select
	id, name, password, role, created, modified
		from todoapp.users where name = $1`
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}
	return u, nil
}
