package main

import ("fmt"
    "net/http"
)

func main(){
  server := http.Server{
    Addr: ":3000",
    Handler: http.HandlerFunc(basicHandler),
  }
  fmt.Println("Server is running on port 3000")
  err := server.ListenAndServe()
  if err != nil {
    fmt.Println(err)
  }
}

func basicHandler(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("Hello World"))
}
