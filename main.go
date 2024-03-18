package main

import ("fmt"
  "context"
"github.com/bun137/microservicesInGo/application"
)

func main(){
  app := application.New()
  ctx := context.Background()
  err := app.Start(ctx)
  if err != nil {
    fmt.Println(err)
  }
}
