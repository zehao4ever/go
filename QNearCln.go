package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func getMsg() (msg string) {
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	msg = strings.Replace(msg, "\n", "", -1)

	return
}
func recv(clnSck net.Conn) {
	// read msg from server.
	buf := make([]byte, 1024)
	for {
		dataLen, err := clnSck.Read(buf)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		var d DataFrame
		if err = json.Unmarshal(buf[:dataLen], &d); err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(d.Msg, "from", d.Src)
	}
}
func cln() {
	// ip, Domain Name, www.baidu.com
	srv, err := net.ResolveTCPAddr("tcp", srvAddr)
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

	// send msg to server.
	transMsg := DataFrame{Src: me}

	var dstPos, msgPos int
	for {
		msg := strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(getMsg(), "　", " ", -1),
						"，", ",", -1),
					"＠", "@", -1),
				", ", ",", -1),
			" ,", ",", -1)

		dstPos = strings.Index(msg, "@")

		if dstPos == 0 {
			msgPos = strings.IndexAny(msg, " ")
			if msgPos <= 0 {
				fmt.Println("empty msg")
				continue
			}
			transMsg.Dst = strings.Split(msg[1:msgPos], ",")
			transMsg.Msg = msg[msgPos:]
		} else {
			transMsg.Msg = msg
		}
		if len(transMsg.Msg) <= 0 {
			fmt.Println("empty msg")
			continue
		}

		if err := transMsg.Marshal(); err != nil {
			fmt.Println(err.Error())
			continue
		}

		clnSck.Write(transMsg.dataBuf)
	}
	clnSck.Close()
}

type DataFrame struct {
	Sign    int      `json:"S"`
	Length  string   `json:"L"`
	Src     string   `json:"src",omitempty`
	Dst     []string `json:"dst",omitempty`
	Msg     string   `json:"msg",omitempty`
	dataBuf []byte
}

func (d *DataFrame) Marshal() error {
	if len(d.Src) <= 0 || d.Dst == nil || len(d.Dst) <= 0 || len(d.Msg) <= 0 {
		return errors.New("param Src/Dst/Msg is empty!")
	}

	d.Sign = 142857
	d.Length = "0x12345678"
	var err error
	d.dataBuf, err = json.Marshal(d)
	if err != nil {
		return err
	}

	d.Length = fmt.Sprintf("%#08X", len(d.dataBuf))
	d.dataBuf, err = json.Marshal(d)

	return err
}

var me string
var srvAddr string

func main() {
	me = "mickey"
	srvAddr = "127.0.0.1:6666"

	if len(os.Args) == 3 {
		srvAddr = os.Args[1] + ":6666"
		me = os.Args[2]
	} else if len(os.Args) == 2 {
		me = os.Args[1]
	}
	cln()
}
