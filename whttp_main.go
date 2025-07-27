package main

import (
	"WUtils/WHttp"
	"fmt"
)

func main() {
	err := WHttp.StartServer(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
