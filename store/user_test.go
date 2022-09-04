package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/jmoiron/sqlx"
)

func TestRepository_RegisterUser(t *testing.T) {
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
	var wantID int64 = 30

	type args struct {
		ctx context.Context
		db  Execer
		u   *entity.User
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"ok case ": {
			fields: fields{Clocker: fc},
			args: args{
				ctx: ctx,
				db:  xdb,
				u: &entity.User{
					Name:     "test-user",
					Password: "test",
					Role:     "user",
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
			mock.ExpectQuery(
				regexp.QuoteMeta(`insert into todoapp.users( name, password, role, created, modified )
				values ($1, $2, $3, $4, $5) returning id;`),
			).WithArgs(tt.args.u.Name, tt.args.u.Password,
				tt.args.u.Role, tt.args.u.Created, tt.args.u.Modified).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(wantID))
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			if err := r.RegisterUser(tt.args.ctx, tt.args.db, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Repository.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
