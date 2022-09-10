package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
	"github.com/chmikata/go_todo_app/testutil/fixture"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserId {
	t.Helper()

	u := fixture.User(nil)
	err := db.QueryRowxContext(ctx,
		`insert into todoapp.users(
		name, password, role, created, modified
		) values ($1, $2, $3, $4, $5) returning id;`,
		u.Name, u.Password, u.Role, u.Created, u.Modified,
	).Scan(&u.ID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return entity.UserId(u.ID)
}

func prepareTasks(ctx context.Context, t *testing.T, db Execer) (entity.UserId, entity.Tasks) {
	t.Helper()

	userId := prepareUser(ctx, t, db)
	oterUserId := prepareUser(ctx, t, db)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserId: userId,
			Title:  "want task 1", Stat: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			UserId: userId,
			Title:  "want task 2", Stat: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserId: oterUserId,
			Title:  "not want task", Stat: "doto",
			Created: c.Now(), Modified: c.Now(),
		},
		wants[1],
	}
	rows, err := db.QueryxContext(ctx,
		`insert into todoapp.tasks
		(user_id, title, stat, created, modified) values
		($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10), ($11, $12, $13, $14, $15) returning id;`,
		tasks[0].UserId, tasks[0].Title, tasks[0].Stat, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserId, tasks[1].Title, tasks[1].Stat, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserId, tasks[2].Title, tasks[2].Stat, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(&tasks[i].ID); err != nil {
			t.Fatal(err)
		}
	}
	return userId, wants
}

func TestRepository_ListTasks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wantUserId, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserId)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got + want)\n%s", d)
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
	tests := map[string]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		"ok case": {
			fields: fields{Clocker: fc},
			args: args{
				ctx: ctx,
				db:  xdb,
				t: &entity.Task{
					UserId:   10,
					Title:    "ok test",
					Stat:     "todo",
					Created:  fc.Now(),
					Modified: fc.Now(),
				},
			},
			wantErr: false,
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			mock.ExpectQuery(
				regexp.QuoteMeta(`insert into todoapp.tasks (user_id, title, stat, created, modified)
				values($1, $2, $3, $4, $5) returning id;`),
			).WithArgs(tt.args.t.UserId, tt.args.t.Title, tt.args.t.Stat, tt.args.t.Created, tt.args.t.Modified).
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
