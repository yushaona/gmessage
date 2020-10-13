package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/yushaona/gjson"
)

func main() {
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	var json gjson.GJSON
	json.SetInt("funcid", 1)
	json.SetString("userid", "123")
	json.SetString("passwd", "123456")
	conn.Write([]byte(json.ToString()))

	time.Sleep(1 * time.Second)
	var d gjson.GJSON
	d.SetString("userid", "123")
	d.SetInt("funcid", 10)
	conn.Write([]byte(d.ToString()))

	go keeplive(conn)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil || err == io.EOF {
			fmt.Printf("%v \n", err)
			break
		}

		fmt.Printf("%s \n", string(buf[:n]))
	}
}

func keeplive(c net.Conn) {

	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			var json gjson.GJSON
			json.SetInt("funcid", 2)
			c.Write([]byte(json.ToString()))
		}
	}
}
