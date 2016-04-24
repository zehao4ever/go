// cstSrv

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func setupDB() {
	var err error
	rootDbPwd := "abc123"
	connStr := "root:" + rootDbPwd + "@/mysql?charset=utf8&loc=Local&parseTime=true"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	cr_db := "CREATE DATABASE IF NOT EXISTS qnearBE DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"
	stmt, err := db.Prepare(cr_db)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt.Close()
	grantSQL := "grant all on qnearBE.* to cstAdmin identified by 'cstDb4ever';"
	stmt, err = db.Prepare(grantSQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()
	grantSQL = "grant all on qnearBE.* to cstAdmin@'' identified by 'cstDb4ever';"
	stmt, err = db.Prepare(grantSQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	grantSQL = "grant all on qnearBE.* to cstAdmin@'localhost' identified by 'cstDb4ever';"
	stmt, err = db.Prepare(grantSQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	grantSQL = "grant all on qnearBE.* to cstAdmin@'127.0.0.1' identified by 'cstDb4ever';"
	stmt, err = db.Prepare(grantSQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()
	db.Close()

	dbPwd := "cstDb4ever"
	connStr = "cstAdmin:" + dbPwd + "@/qnearBE?charset=utf8&loc=Local&parseTime=true"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	cr_table := "create table if not exists t_msg(msg_id int auto_increment primary key, peer varchar(64),msg varchar(128),recvTime datetime not null default 0)"
	stmt, err = db.Prepare(cr_table)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func keepMsg(arg_msg string, arg_peer string) {
	sql := "insert into t_msg(peer,msg) values(?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(arg_peer, arg_msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func queryMsg(arg_peer string) string {
	sql := "select msg,recvTime from t_msg"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}
	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	msg := ""
	recvTime := time.Now()
	msgList := ""
	for rows.Next() {
		rows.Scan(&msg, &recvTime)
		msgList = msgList + recvTime.Format("15:04:05 ") + msg + "\r\n"
	}
	return msgList
}

var clnOnLineChannel chan net.Conn
var clnOffLineChannel chan net.Conn
var msgChannel chan string

func showOnLines(arg_conns map[string]net.Conn) {
	fmt.Println("Online Number: " + strconv.Itoa(len(arg_conns)))
}

func clnMgr() {
	clnOnLineChannel = make(chan net.Conn)
	clnOffLineChannel = make(chan net.Conn)
	msgChannel = make(chan string)

	connList := make(map[string]net.Conn)

	for {
		select {
		case clnSck := <-clnOnLineChannel:
			clnSap := clnSck.RemoteAddr().String()
			fmt.Println(clnSap + " online")
			connList[clnSap] = clnSck
			showOnLines(connList)

		case clnSck := <-clnOffLineChannel:
			clnSap := clnSck.RemoteAddr().String()
			fmt.Println(clnSap + " offline")

			delete(connList, clnSap)
			clnSck.Close()
			showOnLines(connList)

		case msg := <-msgChannel:
			bMsg := []byte(msg)
			for _, v := range connList {
				_, err := v.Write(bMsg)
				if err != nil {
					fmt.Println(err)
					clnOffLineChannel <- v
				}
			}
		}
	}
}

func recv(clnSck net.Conn) {
	buf := make([]byte, 1024)
	for {
		dataLen, err := clnSck.Read(buf)
		if err != nil {
			fmt.Println(err)
			clnOffLineChannel <- clnSck
			return
		}

		msg := string(buf[:dataLen])
		fmt.Println(msg)
		sap := clnSck.RemoteAddr().String()
		keepMsg(msg, sap)

		if len(msg) == 5 && msg[0:5] == "query" {
			msgList := queryMsg(sap)
			clnSck.Write([]byte(msgList))
		}
	}
}

func getMsg() (msg string) {
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func msgToCln() {
	for {
		msg := getMsg()
		msgChannel <- msg
	}
}

func main() {

	setupDB()

	//socket
	srvSck, err := net.Listen("tcp", ":6666")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srvSck.Close()

	go clnMgr()
	go msgToCln()

	for {
		clnSck, err := srvSck.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		clnOnLineChannel <- clnSck
		go recv(clnSck)
	}

}
