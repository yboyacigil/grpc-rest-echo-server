package main

import (
	"fmt"

	"github.com/yboyacigil/grpc-rest-echo-server/server"
)

func main() {
	e := server.New()
	e.Start()
	fmt.Println("EchoServer started with GRPC and REST endpoints.")
	e.WaitStop()
}
