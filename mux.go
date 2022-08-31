package main

import (
	"context"
	"net/http"

	"github.com/chmikata/go_todo_app/clock"
	"github.com/chmikata/go_todo_app/config"
	"github.com/chmikata/go_todo_app/handler"
	"github.com/chmikata/go_todo_app/service"
	"github.com/chmikata/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := &store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServedHTTP)
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: r},
	}
	mux.Get("/tasks", lt.ServedHTTP)
	return mux, cleanup, nil
}
