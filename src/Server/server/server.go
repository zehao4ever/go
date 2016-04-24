package server

import (
	"Server/controller"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	//	"strconv"
)

const (
	ACTION_LOGIN             = 1
	ACTION_REGISTER          = 2
	ACTION_MODIFY_INFO       = 3
	ACTION_MODIFY_PASS       = 4
	ACTION_FIND_USER         = 5
	ACTION_MODIFY_OTHER_INFO = 6
)

const (
	HANDER     = "142857"
	HANDER_LEN = 6
	DATA_LEN   = 4
)

var (
	g_CliDataChannel  chan controller.TranData
	g_CliLoginChannel chan net.Conn
	g_CliOffChannel   chan net.Conn
	g_id_addr         map[int]string
	g_addr_conn       map[string]net.Conn
)

func init() {
	go prepare()
}
func prepare() {
	g_CliDataChannel = make(chan controller.TranData, 200)
	g_CliLoginChannel = make(chan net.Conn)
	g_CliOffChannel = make(chan net.Conn)
	g_id_addr = make(map[int]string)
	g_addr_conn = make(map[string]net.Conn)

	for {
		select {
		case tranData := <-g_CliDataChannel:
			fmt.Println("-----------------")
			fmt.Println(tranData)
			switch tranData.ActionType {
			case ACTION_LOGIN:
				g_id_addr[tranData.SendId] = tranData.Addr
				g_addr_conn[tranData.Addr] = tranData.Conn
				newTran, err := controller.Login(tranData)
				if err != nil {
					log.Println(err)
				}
				buf, err := Pack(newTran)
				if err != nil {
					log.Println(err)
				}
				tranData.Conn.Write(buf)
			case ACTION_REGISTER:
				newTran, err := controller.Register(tranData)
				if err != nil {
					log.Println(err)
				}
				buf, err := Pack(newTran)
				if err != nil {
					log.Println(err)
				}
				tranData.Conn.Write(buf)
			case ACTION_MODIFY_PASS: //修改密码
				newTran, err := controller.ModifyPass(tranData)
				if err != nil {
					log.Println(err)
				}
				buf, err := Pack(newTran)
				if err != nil {
					log.Println(err)
				}
				tranData.Conn.Write(buf)
			case ACTION_MODIFY_INFO:
				newTran, err := controller.ModifyUserInfo(tranData)
				if err != nil {
					log.Println(err)
				}
				buf, err := Pack(newTran)
				if err != nil {
					log.Println(err)
				}
				tranData.Conn.Write(buf)
			case ACTION_MODIFY_OTHER_INFO:
				newTran, err := controller.ModifyOtherUserInfo(tranData)
				if err != nil {
					log.Println(err)
				}
				buf, err := Pack(newTran)
				if err != nil {
					log.Println(err)
				}
				tranData.Conn.Write(buf)
			}

		case conn := <-g_CliLoginChannel:
			g_addr_conn[conn.RemoteAddr().String()] = conn
			fmt.Println(conn.RemoteAddr().String(), " connected")
		case conn := <-g_CliOffChannel:
			delete(g_addr_conn, conn.RemoteAddr().String())
		}
	}

}

func Start() error {
	listen, e := net.Listen("tcp", ":9999")
	if e != nil {
		return e
	}

	for {
		conn, er := listen.Accept()
		if er != nil {
			return er
		}
		go recv(conn)
		g_CliLoginChannel <- conn
	}
}

func recv(arg_conn net.Conn) {
	buf := make([]byte, 1024)
	tempBuf := make([]byte, 0)
	for {
		bufLen, err := arg_conn.Read(buf)
		if err != nil {
			log.Println(err)
			g_CliOffChannel <- arg_conn
			return
		}
		fmt.Println(string(buf))
		buff := append(tempBuf, buf[:bufLen]...)
		fmt.Println(string(buff))
		tempBuf, err = Unpack(buff, arg_conn)
		if err != nil {
			log.Println(err)
		}
	}
}
func Pack(arg_tran controller.TranData) ([]byte, error) {
	dataBuf, err := json.Marshal(arg_tran)
	if err != nil {
		return nil, err
	}
	dataLen := len(dataBuf)
	dataLenBytes := IntToBytes(dataLen)
	buf := append([]byte(HANDER), dataLenBytes...)
	return append(buf, dataBuf...), nil
}

func Unpack(arg_buf []byte, arg_conn net.Conn) ([]byte, error) {
	totalLen := len(arg_buf)
	var i int
	fmt.Println(string(arg_buf[:6]))
	for i = 0; i < totalLen; i++ {
		if totalLen < i+HANDER_LEN+DATA_LEN {
			break
		}
		if string(arg_buf[i:i+HANDER_LEN]) == HANDER {
			dataLen := BytesToInt(arg_buf[i+HANDER_LEN : i+HANDER_LEN+DATA_LEN])
			if totalLen < dataLen+HANDER_LEN+DATA_LEN {
				break
			}
			data := arg_buf[i+HANDER_LEN+DATA_LEN : i+HANDER_LEN+DATA_LEN+dataLen]
			var tranData controller.TranData
			if err := json.Unmarshal(data, &tranData); err != nil {
				log.Println(err)
				return arg_buf[i+HANDER_LEN+DATA_LEN+dataLen:], err
			}
			tranData.Addr = arg_conn.RemoteAddr().String()
			tranData.Conn = arg_conn
			g_CliDataChannel <- tranData
			i = i + HANDER_LEN + DATA_LEN + dataLen - 1
		}
	}
	if i == totalLen {
		return make([]byte, 0), nil
	}
	return arg_buf[i:], nil
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	fmt.Println("ss")
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	fmt.Println("end")
	return int(x)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
