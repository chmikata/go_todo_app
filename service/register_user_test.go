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

func TestRegisterUser_RegisterUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		DB   store.Execer
		Repo UserRegister
	}
	type args struct {
		ctx      context.Context
		name     string
		password string
		role     string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		"ok case": {
			args: args{
				ctx:      context.Background(),
				name:     "user1",
				password: "user1",
				role:     "user",
			},
			want: &entity.User{
				ID:       20,
				Name:     "user1",
				Password: "user1",
				Role:     "user",
				Created:  clock.FixedClocker{}.Now(),
				Modified: clock.FixedClocker{}.Now(),
			},
			wantErr: false,
		},
		"error case": {
			args: args{
				ctx:      context.Background(),
				name:     "user1",
				password: "user1",
				role:     "admin",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			moq := &UserRegisterMock{}
			moq.RegisterUserFunc = func(
				ctx context.Context, db store.Execer, u *entity.User,
			) error {
				if !tt.wantErr {
					fc := clock.FixedClocker{}
					u.ID = tt.want.ID
					u.Created = fc.Now()
					u.Modified = fc.Now()
					return nil
				}
				return errors.New("error from mock")
			}
			ru := &RegisterUser{
				DB:   tt.fields.DB,
				Repo: moq,
			}
			got, err := ru.RegisterUser(tt.args.ctx, tt.args.name, tt.args.password, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// hack because password can not be fixed
			if got != nil {
				got.Password = tt.want.Password
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser.RegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
