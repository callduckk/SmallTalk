package main

import (
	"SmallTalk/internal/server"
	"fmt"
)

func main() {
	server.RunServer("127.0.0.1:9000")
	fmt.Scanf("%s")
}
