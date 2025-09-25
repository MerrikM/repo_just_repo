// echo.go
package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Echo server running on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			for {
				n, err := c.Read(buf)
				if err != nil {
					return
				}
				c.Write(buf[:n])
			}
		}(conn)
	}
}
