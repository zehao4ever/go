// UserManaServer project UserManaServer.go
package main

import (
	"UserManaServer/process"
	// "io/ioutil"
	"fmt"
	"log"
	"net"
)

func main() {
	//设置监听端口9999
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Println("-------------")
		go process.HandleConn(conn)
	}

}
