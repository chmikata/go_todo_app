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
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `insert into todoapp.users(
		name, password, role, created, modified
	) values ($1, $2, $3, $4, $5) returning id;`
	err := db.QueryRowxContext(
		ctx, sql, u.Name, u.Password,
		u.Role, u.Created, u.Modified,
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
