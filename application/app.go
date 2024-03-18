package application

import ("net/http"
"context"
"fmt"
"github.com/redis/go-redis/v9"
  "time"
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
  err:= a.rdb.Ping(ctx).Err()
  if err != nil {
    return fmt.Errorf("redis error: %w", err)
  }
  
  defer func(){
    err := a.rdb.Close()
    if err != nil {
      fmt.Println("redis error: ", err)
    }
  }()

  fmt.Println("redis connected")

  ch := make(chan error, 1)

  go func() {
  err = server.ListenAndServe()
  if err != nil {
    ch <- fmt.Errorf("server error: %w", err)
  }
    close(ch)
 }()

  select{
 case err = <-ch:
    return err
  case <-ctx.Done():
    timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return server.Shutdown(timeout)
}

  return nil
}
