package store

import (
	"context"

	"github.com/chmikata/go_todo_app/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `select
		id, title,
		stat, created, modified
		from todoapp.tasks;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `insert into todoapp.tasks
		(title, stat, created, modified)
		values($1, $2, $3, $4) returning id;`
	err := db.QueryRowxContext(
		ctx, sql, t.Title, t.Stat,
		t.Created, t.Modified,
	).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}
