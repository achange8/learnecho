package main

import (
	"fmt"
	"os"

	"github.com/achange8/learnecho/router"
)

func main() {
	fmt.Println("Welcome osh server with echo!")
	fmt.Println(os.Getenv("DBNAME"))
	e := router.New()
	e.Logger.Fatal(e.Start(":8080"))
}
