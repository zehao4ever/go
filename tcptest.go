package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

var parten string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var op = [2]string{"login", "register"}
var count int = 0
var wg sync.WaitGroup

func getInfo() (name, pass, opr string) {
	//length of username
	nameLen := rand.Intn(16) + 1
	//length of password
	passLen := rand.Intn(16) + 1
	//the string of name and pass
	name = ""
	pass = ""

	for i := 0; i < nameLen; i++ {
		//index of the parten
		index := rand.Intn(len(parten))
		name += string(parten[index])
	}

	for j := 0; j < passLen; j++ {
		//index of the parten
		index := rand.Intn(len(parten))
		pass += string(parten[index])
	}
	//index of the url
	pos := rand.Intn(2)

	return name, pass, op[pos]
}

func test() {
	conn, err := net.Dial("tcp", "localhost:9998")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		count++
		fmt.Println(count)
	}

	//	wg.Add(-1)
	name, pass, op := getInfo()
	//	op = "register"
	var content string = "{\"Op\":\"" + op + "\",\"Username\":\"" + name + "\",\"Password\":\"" + pass + "\"}"
	conn.Write([]byte(content))

	b := make([]byte, 1024)

	for {
		t := rand.Int63n(10)
		time.Sleep(time.Duration(t) * time.Second)
		name, pass, op := getInfo()
		//  op = "register"
		var content string = "{\"Op\":\"" + op + "\",\"Username\":\"" + name + "\",\"Password\":\"" + pass + "\"}"
		conn.Write([]byte(content))

		_, err := conn.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	time.Sleep(60 * time.Second)
	wg.Add(-1)
}

func main() {

	for i := 0; i < 20000; i = i + 1 {
		wg.Add(1)
		go test()
	}

	time.Sleep(30 * time.Second)
	wg.Wait()
}
