package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func getMsg() (msg string) {
	reader := bufio.NewReader(os.Stdin)

	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

func recv(clnSck net.Conn) {
	// read msg from server.
	buf := make([]byte, 1024)
	for {
		dataLen, err := clnSck.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:dataLen]))
	}
}

func cln() {
	// ip, Domain Name, www.baidu.com

	srv, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6666")
	if err != nil {
		fmt.Println(err)
		return
	}

	clnSck, err := net.DialTCP("tcp", nil, srv)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer clnSck.Close()
	go recv(clnSck)

	next <- true
	// send msg to server.
	for {
		msg := getMsg()
		clnSck.Write([]byte(msg))
	}
	clnSck.Close()
}

var next chan bool

func main() {
	next = make(chan bool)
	for {
		go cln()
		<-next
	}
}
