// tcpsrv.go
package main

import (
	"database/sql"
	//  "fmt"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Password string
}

//database operation

var db *sql.DB
var mutex sync.Mutex

//init the connection to the database
func InitDB() (err error) {
	//  var err error
	db, err = sql.Open("mysql", "root:abc123@/cst")
	if err != nil {
		log.Fatal(err)
		return
	}

	db.SetMaxIdleConns(10000) //set the max number of the connection to the database
	db.SetMaxOpenConns(20000)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

//insert a user into the database
func InsertUser(user User) (success bool) {
	stmt, err := db.Prepare("INSERT INTO t_user(t_username,t_password) values(?,?)")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Username, user.Password)
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	if n > 0 {
		success = true
	}
	return
}

//update password of the user
func ModifyPass(user User) (success bool) {
	stmt, err := db.Prepare("UPDATE t_user SET t_password=? WHERE t_username=?")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Password, user.Username)
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	if n > 0 {
		success = true
	}
	return
}
func CloseDB() {
	db.Close()
}

//users operation
//
//
func Login(user User) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := Cache[user.Username]
	return ok
}

func Register(user User) bool {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := Cache[user.Username]
	if ok {
		return false
	} else {
		Cache[user.Username] = user.Password
		intoDB_chan <- user
		return true
	}
}

func ResetPass(user User) bool {
	_, ok := Cache[user.Username]
	if ok {
		Cache[user.Username] = user.Password
		updateDB_chan <- user
		return true
	} else {
		return false
	}
}

// network operation
//
//
var AllConnMap map[net.Addr]net.Conn

const (
	login_cmd      = "login"
	register_cmd   = "register"
	exit_cmd       = "exit"
	reset_pass_cmd = "reset_pass"
	MAX_CONN_NUM   = 100000
)

func handler(conn net.Conn) {
	//  address := conn.RemoteAddr()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			//  delete(AllConnMap, address)

			//  fmt.Println("conn read error : ", err)
			//  fmt.Println(conn.RemoteAddr().String(), "unlink")
			return
		}
		buffer := buf[:n]
		go handeRequest(conn, buffer)
	}
}
func handeRequest(conn net.Conn, buf []byte) {

	Buff := bytes.NewReader(buf)
	d := json.NewDecoder(Buff)

	var msg interface{}
	if err := d.Decode(&msg); err != nil {
		log.Println(conn.RemoteAddr().String(), " connection error: ", err)
		return
	}
	//  log.Println(conn.RemoteAddr().String(), "receive data:", msg)
	data := msg.(map[string]interface{})

	cmd_str := data["Op"].(string)
	var user = User{data["Username"].(string), data["Password"].(string)}

	//	fmt.Println("Op: ", cmd_str, " Username: ", user.Username, " Password: ", user.Password)
	switch cmd_str {
	case login_cmd:
		if Login(user) {
			conn.Write([]byte("success\n"))
		} else {
			conn.Write([]byte("fail\n"))
		}
	case register_cmd:
		if Register(user) {
			conn.Write([]byte("success\n"))
		} else {
			conn.Write([]byte("fail\n"))
		}
	case exit_cmd:
		conn.Close()
	case reset_pass_cmd:
		if ResetPass(user) {
			conn.Write([]byte("success\n"))
		} else {
			conn.Write([]byte("fail\n"))
		}
	}
}

//cache operator
//
//
var Cache map[string]string

func InitCache() {
	Cache = make(map[string]string, 1000000)
	stmt, err := db.Prepare("SELECT * FROM t_user")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	var id int
	var name string
	var pass string

	for rows.Next() {
		rows.Scan(&id, &name, &pass)
		Cache[name] = pass
	}
	return
}

var (
	intoDB_chan   = make(chan User, 10000)
	updateDB_chan = make(chan User, 10000)
)

func main() {
	InitDB()
	defer CloseDB()

	InitCache()

	var cur_conn_num int = 0
	conn_chan := make(chan net.Conn, 10000)
	ch_conn_change := make(chan int)

	go func() {
		for conn_change := range ch_conn_change {
			cur_conn_num += conn_change
		}
	}()

	AllConnMap = make(map[net.Addr]net.Conn)

	ser, err := net.Listen("tcp", ":9998")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ser.Close()

	go func() {
		for _ = range time.Tick(5e9) {
			fmt.Println("now connections number is : ", cur_conn_num)
			fmt.Println("data number is :", len(Cache))
		}
	}()
	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range conn_chan {
				ch_conn_change <- 1
				handler(conn)
				ch_conn_change <- -1
			}
		}()

	}
	for i := 0; i < 10; i++ {
		go func() {
			for user := range intoDB_chan {
				InsertUser(user)
			}
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			for user := range updateDB_chan {
				ModifyPass(user)
			}
		}()
	}

	fmt.Println("running...........")
	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Fatal(err)
			fmt.Println("connection err")
			continue
		}
		//  AllConnMap[conn.RemoteAddr()] = conn
		//  fmt.Println("now conn numbers : ", len(AllConnMap))
		//  go handler(conn)
		conn_chan <- conn
	}
}
