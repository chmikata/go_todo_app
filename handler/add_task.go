package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type AddTask struct {
	Service   AddTaskService
	Validator *validator.Validate
}

func (at *AddTask) ServedHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: int(t.ID)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
