package store

import (
	"context"

	"github.com/chmikata/go_todo_app/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer, id entity.UserId,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `select
		id, user_id, title,
		stat, created, modified
		from todoapp.tasks
		where user_id = $1;`
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	sql := `insert into todoapp.tasks
		(user_id, title, stat, created, modified)
		values($1, $2, $3, $4, $5) returning id;`
	err := db.QueryRowxContext(
		ctx, sql, t.UserId, t.Title, t.Stat,
		r.Clocker.Now(), r.Clocker.Now(),
	).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}
