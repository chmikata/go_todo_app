package store

import (
	"context"
	"reflect"
	"testing"

	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	type fields struct {
		Clocker clock.Clocker
	}
	type args struct {
		ctx context.Context
		db  Queryer
	}
	fc := clock.FixedClocker{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Tasks
		wantErr bool
	}{
		{
			name:   "Dummy DB test",
			fields: fields{Clocker: fc},
			args: args{
				ctx: ctx,
				db:  tx,
			},
			want: entity.Tasks{
				&entity.Task{
					Title: "task1", Stat: "todo",
					Created: fc.Now(), Modified: fc.Now(),
				},
				&entity.Task{
					Title: "task2", Stat: "todo",
					Created: fc.Now(), Modified: fc.Now(),
				},
				&entity.Task{
					Title: "task3", Stat: "done",
					Created: fc.Now(), Modified: fc.Now(),
				},
			},
			wantErr: false,
		},
	}
	// initialize test case
	if _, err := tx.ExecContext(ctx, "delete from todoapp.tasks;"); err != nil {
		t.Errorf("failed to initialize tasks: %v", err)
	}
	wants := tests[0].want
	rows, err := tx.QueryContext(ctx,
		`insert into todoapp.tasks (title, stat, created, modified) values
		($1, $2, $3, $4),($5, $6, $7, $8), ($9, $10, $11, $12) returning id;`,
		wants[0].Title, wants[0].Stat, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Stat, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Stat, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; rows.Next(); i++ {
		rows.Scan(&wants[i].ID)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			got, err := r.ListTasks(tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// PostgreSQLのtimestampとのズレを解消
			for _, g := range got {
				g.Created = g.Created.UTC()
				g.Modified = g.Modified.UTC()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.ListTasks() = %v, want %v", got[0], tt.want[0])
			}
		})
	}
}

func TestRepository_AddTask(t *testing.T) {
	type fields struct {
		Clocker clock.Clocker
	}
	type args struct {
		ctx context.Context
		db  Execer
		t   *entity.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			if err := r.AddTask(tt.args.ctx, tt.args.db, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Repository.AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
