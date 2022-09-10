package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/chmikata/go_todo_app/auth"
	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/store"
)

func TestAddTask_AddTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		DB   store.Execer
		Repo TaskAdder
	}
	type args struct {
		ctx   context.Context
		title string
	}

	tests := map[string]struct {
		fields  fields
		args    args
		want    *entity.Task
		wantErr bool
	}{
		"ok case": {
			args: args{ctx: auth.SetUserID(context.Background(), 10), title: "test title"},
			want: &entity.Task{
				ID:       10,
				UserId:   10,
				Title:    "test title",
				Stat:     entity.TaskStatusTodo,
				Created:  clock.FixedClocker{}.Now(),
				Modified: clock.FixedClocker{}.Now(),
			},
			wantErr: false,
		},
		"error case": {
			args:    args{context.Background(), "test title"},
			want:    nil,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			moq := &TaskAdderMock{}
			moq.AddTaskFunc = func(
				ctx context.Context, db store.Execer, t *entity.Task,
			) error {
				if !tt.wantErr {
					fc := clock.FixedClocker{}
					t.ID = tt.want.ID
					t.UserId = tt.want.UserId
					t.Created = fc.Now()
					t.Modified = fc.Now()
					return nil
				}
				return errors.New("error from mock")
			}
			at := &AddTask{
				DB:   tt.fields.DB,
				Repo: moq,
			}
			got, err := at.AddTask(tt.args.ctx, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTask.AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddTask.AddTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
