package service

import (
	"context"

	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskLister TaskAdder UserRegister UserGetter TokenGenerator
type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, id entity.UserId) (entity.Tasks, error)
}

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}
