package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type DataFrame struct {
	Sign    int    `json:"S"`
	Length  string `json:"L"`
	Op      string `json:"op"`
	UserId  int    `json:"user_id"`
	Pwd     string `json:"pwd"`
	Result  int    `json:"res"`
	DataBuf []byte
	Sap     string
}

func (this *DataFrame) Marshal() error {
	if len(this.Op) <= 0 || this.UserId == 0 {
		return errors.New("I don't kown what is the op or user")
	}
	this.Sign = 142857
	this.Length = "0x12345678"
	var err error
	this.DataBuf, err = json.Marshal(this)
	if err != nil {
		return err
	}
	this.Length = fmt.Sprintf("%#08x", len(this.DataBuf))
	this.DataBuf, err = json.Marshal(this)
	return err
}

func recv(arg_conn net.Conn) {

}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("输入有误v，请输入用户名和密码")
		return
	}

}
