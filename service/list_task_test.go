package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/store"
)

func TestListTask_ListTasks(t *testing.T) {
	t.Parallel()

	type fields struct {
		DB   store.Queryer
		Repo TaskLister
	}
	type args struct {
		ctx context.Context
	}
	fc := clock.FixedClocker{}
	tests := map[string]struct {
		fields  fields
		args    args
		want    entity.Tasks
		wantErr bool
	}{
		"ok case1": {
			args: args{context.Background()},
			want: entity.Tasks{
				&entity.Task{
					ID:       10,
					Title:    "task1",
					Stat:     entity.TaskStatusTodo,
					Created:  fc.Now(),
					Modified: fc.Now(),
				},
				&entity.Task{
					ID:       11,
					Title:    "task2",
					Stat:     entity.TaskStatusTodo,
					Created:  fc.Now(),
					Modified: fc.Now(),
				},
				&entity.Task{
					ID:       12,
					Title:    "task3",
					Stat:     entity.TaskStatusTodo,
					Created:  fc.Now(),
					Modified: fc.Now(),
				},
			},
			wantErr: false,
		},
		"ok case2": {
			args:    args{context.Background()},
			want:    entity.Tasks{},
			wantErr: false,
		},
		"error case": {
			args:    args{context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			moq := &TaskListerMock{}
			moq.ListTasksFunc = func(
				ctx context.Context, db store.Queryer,
			) (entity.Tasks, error) {
				if !tt.wantErr {
					return tt.want, nil
				}
				return nil, errors.New("error from mock")
			}
			lt := &ListTask{
				DB:   tt.fields.DB,
				Repo: moq,
			}
			got, err := lt.ListTasks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTask.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTask.ListTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
