package process

import (
	"UserManaServer/control"
	// "io/ioutil"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	END_MARK = "#" //数据的结束标志
)

var (
	//全局变量以用户id作为key，token作为value
	g_nameId_token map[int]string = make(map[int]string)
	g_token_nameId map[string]int = make(map[string]int)
)

/**
*对客户端的连接进行处理，数据的读取
 */
func HandleConn(arg_conn net.Conn) {
	buf := make([]byte, 1024)
	preDataStr := ""
	for {
		//读取数据
		dataLen, err := arg_conn.Read(buf)
		if err != nil {
			//读取结束，说明客户端连接已断开
			if err == io.EOF {
				log.Println("the conn is closed")
				arg_conn.Close()
			} else {
				log.Println("Read error: ", err)
			}
			break
		}
		buf = buf[:dataLen]
		preDataStr = spiltStringData(arg_conn, preDataStr, buf)
	}
}

/**
*对数据以#作为结束标志进行切割，并返回多余的字符串数据
*参数：第一个是上次切割时多余的字符串，第二个是从客户端刚
*获取的数据
 */
func spiltStringData(arg_conn net.Conn, arg_old_str string, arg_new_data []byte) string {
	newStr := string(arg_new_data)                        //转换为字符串
	strDatasAfterSpilt := strings.Split(newStr, END_MARK) //[]string
	strDatasAfterSpilt[0] = arg_old_str + strDatasAfterSpilt[0]
	arrayLen := len(strDatasAfterSpilt)
	//如果数组长度为1，则数据中没有结束标志，暂存数据，不应处理
	//如果数组长度>=2，分为两种情况：
	//1.若数组的最后一个元素为空（“”），说明数据完整，都应进行处理
	//2.如数组的最后一个元素不为空，说明最后一个数据读取中断，应处理len-1次
	if arrayLen == 1 {
		return strDatasAfterSpilt[0]
	} else {
		var visitLen int = 0
		if strDatasAfterSpilt[arrayLen-1] == "" {
			visitLen = arrayLen - 1
		} else {
			visitLen = arrayLen - 2
		}
		for i := 0; i < visitLen; i++ {
			go handleData(arg_conn, strDatasAfterSpilt[i])
		}
	}
	return strDatasAfterSpilt[arrayLen-1]
}

func handleData(arg_conn net.Conn, arg_data string) {
	//首先要把传进来string类型的数据装换成json
	//然后再把json中的数据放进map[string]string中
	var v interface{}
	err := json.Unmarshal([]byte(arg_data), &v)
	if err != nil {
		log.Println(err)
		return
	}
	data := v.(map[string]interface{})

	//把interface{} 转成string
	strData := make(map[string]string)
	for key := range data {
		strData[key] = data[key].(string)
	}
	op := data[OPERATION].(string) //取出操作命令
	switch op {

	//检测是否为登录请求，是的话不用验证token
	case control.LOGINON:
		//根据op路由到对应的func
		result := control.G_op_func[op](strData)
		success := result[TAG]
		if success == control.LOGIN_SUCCESS {
			userId, _ := strconv.Atoi(strData[USER_ID])
			token := NewToken(userId, strData[PASSWOR])
			result[TOKEN] = token
			g_nameId_token[userId] = token
			g_token_nameId[token] = userId
		}
		Write2Client(arg_conn, result)
	//判断是不是查看或修改自己的信息
	case control.VIEW_ITSELF_INFO, control.MODIFY_ITSELF_INFO:
		//验证token
		userId, _ := strconv.Atoi(strData[USER_ID])
		if CheckToken(userId, strData[TOKEN]) {
			//检测是不是自己对自己的信息进行操作

			//根据op路由到对应的func
			result := control.G_op_func[op](strData)
			Write2Client(arg_conn, result)
		} else {
			Write2ClientWhenTokenOverdue(arg_conn)
		}
	default:
		//验证token
		userId, _ := strconv.Atoi(strData[USER_ID])
		if CheckToken(userId, strData[TOKEN]) {
			//根据op路由到对应的func
			result := control.G_op_func[op](strData)
			Write2Client(arg_conn, result)
		} else {
			Write2ClientWhenTokenOverdue(arg_conn)
		}
	}
}

/**
*根据时间，用户名，密码生成一个token,长度为20
 */
func NewToken(arg_nameId int, arg_pass string) string {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10)+strconv.Itoa(arg_nameId)+arg_pass)
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token[:20]
}

//检测token的正确性
//需要同时验证用户名和对应的token
func CheckToken(arg_nameId int, arg_token string) bool {
	if len(arg_token) > 20 {
		return false
	}
	if g_nameId_token[arg_nameId] == arg_token && g_token_nameId[arg_token] == arg_nameId {
		return true
	}
	return false
}

//把map[string]string转化为json,在末尾加上结束标志
func Map2json(arg_data map[string]string) []byte {
	b, _ := json.Marshal(arg_data)
	str := string(b)
	b = []byte(str + END_MARK)
	return b
}

//写数据到客户端
func Write2Client(arg_conn net.Conn, arg_data map[string]string) {
	dataBuf := Map2json(arg_data)
	arg_conn.Write(dataBuf)
}

//验证token失败后，回应客户端
func Write2ClientWhenTokenOverdue(arg_conn net.Conn) {
	data := make(map[string]string)
	data[TAG] = control.OVERDUE_TOKEN
	Write2Client(arg_conn, data)
}
