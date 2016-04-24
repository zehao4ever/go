package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//数据库连接池
var g_db *sql.DB

const (
	INSERT_USER_SQL     = "INSERT INTO user(name,pwd,sex,age,intro,role_id) VALUES(?,?,?,?,?,?)"
	DEL_USER_SQL        = "DELETE FROM user WHERE user_id = ?"
	QUERY_USER_INFO_SQL = "SELECT *FROM user WHERE user_id = ?"
	MODIFY_USER_sQL     = "UPDATE user SET name=?,sex=?,age=?,intro=?,role_id=? WHERE user_id=?"
)

//导包的时候会自动调用，不需要显式调用init函数，不过数据库需要自己去建立
func init() {
	var err error
	g_db, err = sql.Open("mysql", "root:abc123@/cstdb")
	if err != nil {
		panic(err)
	}
	//set the max number of the connection to the database
	g_db.SetMaxIdleConns(1000)
	g_db.SetMaxOpenConns(2000)
	err = g_db.Ping()
	if err != nil {
		panic(err)
	}
	SQL := `CREATE TABLE IF NOT EXISTS user
      (
      user_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      name varchar(20),
      pwd varchar(20),
      sex int,
      age int,
      intro varchar(100),
      role_id int NOT NULL
      );`
	stmt, err := g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS blacklist
      (
      user_id int NOT NULL PRIMARY KEY,
      join_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
      );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS friendship
      (
      first_user_id int NOT NULL,
      second_user_id int NOT NULL,
      create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (first_user_id,second_user_id),
      FOREIGN KEY (first_user_id) REFERENCES user(user_id) ON DELETE CASCADE,
      FOREIGN KEY (second_user_id) REFERENCES user(user_id) ON DELETE CASCADE
      );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS billboard
       (
       bill_id int PRIMARY KEY AUTO_INCREMENT,
       content varchar(500),
       create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       user_id int,
       user_name varchar(30),
       FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
       );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS role
      (
      role_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      role_name varchar(20) NOT NULL,
      role_desc varchar(100)
      );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS auth
      (
      auth_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      auth_name varchar(20) NOT NULL,
      auth_desc varchar(50)
      );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

	SQL = `CREATE TABLE IF NOT EXISTS role_auth
      (
      role_id int NOT NULL,
      auth_id int NOT NULL,
      PRIMARY KEY (role_id,auth_id)
      );`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()
	//修改表billboard的外键，加入级联删除
	SQL = `ALTER TABLE billboard
			DROP FOREIGN KEY billboard_ibfk_1`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()
	SQL = `ALTER TABLE billboard
			ADD CONSTRAINT billboard_ibfk_1
			FOREIGN KEY (user_id)
			REFERENCES user(user_id) ON DELETE CASCADE`
	stmt, err = g_db.Prepare(SQL)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt.Close()

}

/*
//用户表
mysql> CREATE TABLE user
      (
      user_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      name varchar(20),
      pwd varchar(20),
      sex int,
      age int,
      intro varchar(100),
      role_id int NOT NULL
      );
//黑名单表
mysql> CREATE TABLE blacklist
      (
      user_id int NOT NULL PRIMARY KEY,
      join_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
      );

//好友关系
mysql> CREATE TABLE friendship
      (
      first_user_id int NOT NULL,
      second_user_id int NOT NULL,
      create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (first_user_id,second_user_id),
      FOREIGN KEY (first_user_id) REFERENCES user(user_id) ON DELETE CASCADE,
      FOREIGN KEY (second_user_id) REFERENCES user(user_id) ON DELETE CASCADE
      );

公告
mysql>  CREATE TABLE billboard
       (
       bill_id int PRIMARY KEY AUTO_INCREMENT,
       content varchar(500),
       create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       user_id int,
       user_name varchar(30),
       FOREIGN KEY (user_id) REFERENCES user(user_id)
       );

角色
mysql> CREATE TABLE role
      (
      role_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      role_name varchar(20) NOT NULL,
      role_desc varchar(100)
      );

权限
mysql> CREATE TABLE auth
      (
      auth_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      auth_name varchar(20) NOT NULL,
      auth_desc varchar(50)
      );

角色-权限
mysql> CREATE TABLE role_auth
      (
      role_id int NOT NULL,
      auth_id int NOT NULL,
      PRIMARY KEY (role_id,auth_id)
      );
*/

type User struct {
	Uid          int      `json:"id"`
	Name         string   `json:"userName"`
	Pass         string   `json:"password"`
	Sex          int      `json:"gender"`
	IsOnline     bool     `json:"isOnline"`
	FriendIds    []int    `json:"friendIds"`
	FriendNIames []string `json:""friendNames`
	Auths        []int    `json"auths"`
	Role         int      `json:"role"`
	Age          int      `json:"age"`
	Intro        string   `json:"intro"`
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
