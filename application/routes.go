package application

import (
  "net/http"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/bun137/microservicesInGo/handlers"
)

func loadRoutes() *chi.Mux {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  })
  router.Route("/orders", loadOrderRoutes)
  return r
}

func loadOrderRoutes(router chi.Router) {
  router.Post("/", orderHandler.Create)
  router.Get("/", orderHandler.List)
  router.Get("/{id}", orderHandler.GetByID)
  router.Put("/{id}", orderHandler.UpdateByID)
  router.Delete("/{id}", orderHandler.DeleteByID)
}
