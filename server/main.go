package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yushaona/gmessage/server/config"
	"github.com/yushaona/gmessage/server/handle"
	_ "github.com/yushaona/gmessage/server/models"
)

//
func main() {
	server, err := net.Listen("tcp", config.TCP_HOST)
	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}
	defer server.Close()

	handle.Accept(server)
}
