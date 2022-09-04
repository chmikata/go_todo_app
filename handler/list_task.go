package handler

import (
	"net/http"

	"github.com/chmikata/go_todo_app/entity"
)

type ListTask struct {
	Service ListTasksService
}

type ShowTask struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServedHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := make([]ShowTask, len(tasks))
	for i, t := range tasks {
		rsp[i] = ShowTask{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Stat,
		}
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
