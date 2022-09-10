package main

import (
	"context"
	"net/http"

	"github.com/chmikata/go_todo_app/auth"
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
	var c clock.Clocker = clock.RealClocker{}
	r := store.Repository{Clocker: c}
	rcli, err := store.NewKvs(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, c)
	if err != nil {
		return nil, cleanup, err
	}
	l := &handler.Login{
		Service: &service.Login{
			DB:    db,
			Repo:  &r,
			JWTer: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServedHTTP)

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServedHTTP)

	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Get("/", lt.ServedHTTP)
		r.Post("/", at.ServedHTTP)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})
	return mux, cleanup, nil
}
