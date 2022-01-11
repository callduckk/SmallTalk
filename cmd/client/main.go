package main

import (
	"SmallTalk/internal/client"
	"fmt"
	"log"
)

func main() {
	fmt.Printf("Enter a username: ")
	var name string
	fmt.Scan(&name)
	cl, err := client.NewClient(name, "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	var buffer string
	for {
		fmt.Scanln(&buffer)
		cl.Send(buffer)
	}
}
