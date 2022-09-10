package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
	"github.com/go-playground/validator"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok case": {
			reqFile: "testdata/register_user/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/register_user/ok_rsp.json.golden",
			},
		},
		"badrequest case": {
			reqFile: "testdata/register_user/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/register_user/bad_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/register",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			fc := clock.FixedClocker{}
			moq := &RegisterUserServiceMock{}
			moq.RegisterUserFunc = func(
				ctx context.Context, name, password, role string,
			) (*entity.User, error) {
				if tt.want.status == http.StatusOK {
					return &entity.User{
						ID:       10,
						Name:     "john",
						Password: "test",
						Role:     "user",
						Created:  fc.Now(),
						Modified: fc.Now(),
					}, nil
				}
				return nil, errors.New("error from mock")
			}
			ru := &RegisterUser{
				Service:   moq,
				Validator: validator.New(),
			}
			ru.ServedHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(
				t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
