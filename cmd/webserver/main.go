package main

import (
	"fmt"

	"brianhang.me/facegraph/internal/db"
)

func main() {
	db.Init()

	fmt.Println("Hello world!")
}
