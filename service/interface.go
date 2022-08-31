package service

import (
	"context"

	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder Tasklister
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
