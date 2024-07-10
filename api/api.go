package api

import (
	"log"
	"net/http"

	"github.com/IvanMishnev/go_final_project/handlers"
	"github.com/IvanMishnev/go_final_project/internal/constants"
	"github.com/IvanMishnev/go_final_project/middleware"
	"github.com/go-chi/chi/v5"
)

func StartServer() {
	addr := ":" + constants.Port

	r := chi.NewRouter()
	r.Get("/*", handlers.FileServer)
	r.Get("/api/nextdate", handlers.GetNextDate)

	r.Post("/api/signin", handlers.SignIn)

	r.Route("/api/task", func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/", handlers.GetTask)
		r.Post("/", handlers.AddTask)
		r.Put("/", handlers.EditTask)
		r.Delete("/", handlers.DeleteTask)
		r.Post("/done", handlers.DoneTask)
	})
	r.With(middleware.Auth).Get("/api/tasks", handlers.GetTasks)

	log.Printf("server is listening on: %s", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
