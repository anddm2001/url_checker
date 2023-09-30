package main

import (
	"fmt"
	"url_checker/internal/app"
)

func main() {
	fmt.Println("Init app url checker")

	a := app.New()

	fmt.Println("Success init app url checker")

	a.Run()
}
