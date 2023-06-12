package main

import (
	"fmt"
	"org.example/hello/src/web"
)

func main() {
	err := web.InitAndConfigureRouter().Run(":8080")
	if err != nil {
		fmt.Printf("failed to start server: %s", err.Error())
		return
	}
}
