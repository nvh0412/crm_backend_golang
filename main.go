package main

import (
	"fmt"
	"os"
)

func main() {
  a := App{}

  a.Initializer(
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"),
    os.Getenv("APP_DB_CONNECTION_STRING"),
  )

  fmt.Println("Starting HTTP server at 3000...")

  a.Run(":3000")
}
