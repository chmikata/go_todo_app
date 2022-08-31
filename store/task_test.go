package store

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
	"github.com/jmoiron/sqlx"
)

func TestRepository_ListTasks(t *testing.T) {
	t.Parallel()

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
	// initialize test case. test record delete
	if _, err := tx.ExecContext(ctx, "delete from todoapp.tasks;"); err != nil {
		t.Errorf("failed to initialize todoapp.tasks: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := tx.QueryContext(ctx,
				`insert into todoapp.tasks (title, stat, created, modified) values
				($1, $2, $3, $4),($5, $6, $7, $8), ($9, $10, $11, $12) returning id;`,
				tt.want[0].Title, tt.want[0].Stat, tt.want[0].Created, tt.want[0].Modified,
				tt.want[1].Title, tt.want[1].Stat, tt.want[1].Created, tt.want[1].Modified,
				tt.want[2].Title, tt.want[2].Stat, tt.want[2].Created, tt.want[2].Modified,
			)
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; rows.Next(); i++ {
				if err := rows.Scan(&tt.want[i].ID); err != nil {
					t.Fatal(err)
				}
			}
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			got, err := r.ListTasks(tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Eliminate discrepancies with PostgreSQL timestamp
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
	t.Parallel()

	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	type fields struct {
		Clocker clock.Clocker
	}
	xdb := sqlx.NewDb(db, "postgres")
	fc := clock.FixedClocker{}
	var wantID int64 = 20

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
		{
			name:   "test1",
			fields: fields{Clocker: fc},
			args: args{
				ctx: ctx,
				db:  xdb,
				t: &entity.Task{
					Title:    "ok test",
					Stat:     "todo",
					Created:  fc.Now(),
					Modified: fc.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(
				regexp.QuoteMeta(`insert into todoapp.tasks (title, stat, created, modified)
				values($1, $2, $3, $4) returning id;`),
			).WithArgs(tt.args.t.Title, tt.args.t.Stat, tt.args.t.Created, tt.args.t.Modified).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(wantID))
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			if err := r.AddTask(tt.args.ctx, tt.args.db, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Repository.AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
