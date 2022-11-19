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
  )

  fmt.Println("Start server at 3000...")

  a.Run(":3000")
}
