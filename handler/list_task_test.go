package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
)

func TestListTask(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok case": {
			tasks: []*entity.Task{
				{
					ID:    1,
					Title: "test1",
					Stat:  entity.TaskStatusTodo,
				},
				{
					ID:    2,
					Title: "test2",
					Stat:  entity.TaskStatusDone,
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json.golden",
			},
		},
		"empty case": {
			tasks: []*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_rsp.json.golden",
			},
		},
		"error case": {
			tasks: nil,
			want: want{
				status:  http.StatusInternalServerError,
				rspFile: "testdata/list_task/bad_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListTasksServiceMock{}
			moq.ListTasksFunc = func(ctx context.Context) (entity.Tasks, error) {
				if tt.tasks != nil {
					return tt.tasks, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListTask{Service: moq}
			sut.ServedHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(
				t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
