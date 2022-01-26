package main

import (
	"fmt"

	"github.com/achange8/learnecho/router"
)

func main() {
	fmt.Println("Welcome osh server with echo!")
	e := router.New()

	e.Logger.Fatal(e.Start(":8080"))
}
