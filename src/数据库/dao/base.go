package dao

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//数据库连接池
var g_db *sql.DB

const (
	INSERT_USER_SQL = "INSERT INTO user(name,pwd,sex,age,intro,role_id) VALUES(?,?,?,?,?,?)"
	DEL_USER_SQL    = "DELETE FROM user WHERE user_id = ?"
)

/**
* 初始化数据库，使用数据库接口之前需调用
 */
func InitDB() (err error) {
	g_db, err = sql.Open("mysql", "root:abc123@/cstdb")
	if err != nil {
		log.Println(err)
		return
	}
	//set the max number of the connection to the database
	g_db.SetMaxIdleConns(1000)
	g_db.SetMaxOpenConns(2000)
	err = g_db.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

type User struct {
	Uid   int    `json:"id,string"`
	Name  string `json:"name"`
	Pass  string `json:"pwd"`
	Age   int    `json:"age,string"`
	Sex   int    `json:"sex,string"`
	Intro string `json:"intro"`
	Role  int    `json:"role,string"`
}

type Role struct {
	Rid   int    `json:"role,string"`
	Name  string `json:"role_name"`
	Desc  string `json:"role_desc"`
	Auths []int  `json:"auths"`
}

type Auth struct {
	Aid     int    `json:"auth,string"`
	Name    string `json:"auth_name"`
	Desc    string `json:"auth_desc"`
	Role_id []int  `json:"role"`
}

type Billboard struct {
	Bid        int    `json:"id,string"`
	Content    string `json:"content"`
	CreateTime string `json:"bill_create_time"`
	User_id    int    `json:"user_id,string"`
	User_name  string `json:"name"`
}

// type Billboard struct {
// 	Bid        int    `json:"bill_id,string"`
// 	Title      string `json:"bill_title"`
// 	Cotent     string `json:"bill_content"`
// 	CreateTime string `json:"bill_create_time"`
// 	LastTime   string `json:"bill_last_time"`
// 	User_id    int    `json:"user_id,string"`
// }
