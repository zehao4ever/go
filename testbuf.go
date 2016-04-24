package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	//	"encoding/json"
	"io"
	//	"log"
	//	"net"
	"strconv"
	//	"strings"
	"time"
)

func spiltStringData(arg_old_str string, arg_new_data []byte) string {
	newStr := string(arg_new_data)                   //转换为字符串
	strDatasAfterSpilt := strings.Split(newStr, "#") //[]string
	fmt.Println("len: ", len(strDatasAfterSpilt))
	strDatasAfterSpilt[0] = arg_old_str + strDatasAfterSpilt[0]
	for i := 0; i < len(strDatasAfterSpilt); i++ {
		fmt.Println(strDatasAfterSpilt[i])
	}
	//o handleData(strDatasAfterSpilt[0])
	return ""
}

var g_tokens map[string]string = make(map[string]string)

func handleData(arg_data string) {
	//首先要把传进来string类型的数据装换成json
	//然后再把json中的数据放进map[string]string中
	var v interface{}
	err := json.Unmarshal([]byte(arg_data), &v)
	if err != nil {
		log.Fatal(err)
		return
	}
	//data := v.(map[string]interface{})
	fmt.Println(v)
}
func NewToken(arg_name, arg_pass string) string {
	crutime := time.Now().Unix()
	fmt.Println("crutime-->", crutime)

	h := md5.New()
	fmt.Println("h-->", h)

	fmt.Println("strconv.FormatInt(crutime, 10)-->", strconv.FormatInt(crutime, 10))
	io.WriteString(h, strconv.FormatInt(crutime, 10))

	fmt.Println("h-->", h)

	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("token--->", token)

	fmt.Println(len("8e1a188743c6077110da3c9778183031"))
	return ""
}

//把map[string]string转化为json
func Map2json(arg_data map[string]string) []byte {
	b, _ := json.Marshal(arg_data)
	str := string(b)
	fmt.Println(str)
	b = []byte(str + "#")
	return b
}

func main() {
	// str := []byte("000wwww")
	// spiltStringData("1", str)
	// g_tokens["nnn"] = "55555"
	op := "ppp"
	name := "ppppp"
	// pass := "wwwwwwwwwwww"
	// handleData("{\"Op\":\"" + op + "\",\"Username\":\"" + name + "\",\"Password\":\"" + pass + "\"}")
	//	NewToken(op, name)
	mymap := make(map[string]string)
	mymap[op] = name
	buf := Map2json(mymap)
	fmt.Println(string(buf))
}
