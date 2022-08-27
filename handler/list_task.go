package handler

import (
	"net/http"

	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/store"
)

type ListTask struct {
	Store *store.TaskStore
}

type ShowTask struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServedHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks := lt.Store.All()
	rsp := make([]ShowTask, len(tasks))
	for i, t := range tasks {
		rsp[i] = ShowTask{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		}
	}
	RespondJSON(ctx, w, tasks, http.StatusOK)
}
