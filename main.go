package main

import ("fmt"
  "context"
  "os"
  "os/signal"
"github.com/bun137/microservicesInGo/application"
)

func main(){
  app := application.New()
  ctx,cancel := signal.NotifyContext(context.Background(), os.Interrupt)
  defer cancel()

  err := app.Start(ctx)
  if err != nil {
    fmt.Println(err)
  }
}
