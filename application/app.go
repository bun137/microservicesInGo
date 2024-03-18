package application

import ("net/http"
"context"
"fmt"
"github.com/redis/redis-go/v9"
)

type App struct {
  router http.Handler
  rdb *redis.Client
}

func New() *App {
  app := &App{
    router : loadRoutes(),
    rdb : redis.NewClient(&redis.Options{}),
  }
  return app
}

func (a *App) Start(ctx context.Context) error {
  server := &http.Server{
    Addr: ":3000",
    Handler: a.router,
  }
  err := server.ListenAndServe()
  if err != nil {
    return fmt.Errorf("server error: %w", err)
  }
  return nil
}
