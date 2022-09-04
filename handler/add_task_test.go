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

func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badrequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
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
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			fc := clock.FixedClocker{}
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(
				ctx context.Context, title string,
			) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{
						ID:       1,
						Title:    "Implement a handler",
						Stat:     entity.TaskStatusTodo,
						Created:  fc.Now(),
						Modified: fc.Now(),
					}, nil
				}
				return nil, errors.New("error from mock")
			}

			at := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			at.ServedHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(
				t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
