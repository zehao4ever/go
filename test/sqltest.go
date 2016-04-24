package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB

//init the connection to the database
func InitDB() (err error) {
	//  var err error
	db, err = sql.Open("mysql", "root:abc123@/test")
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

func main() {
	InitDB()
	stmt, err := db.Prepare("select * from t1")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	row, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	var i int
	var s string
	for row.Next() {
		row.Scan(&i, &s)
	}
	fmt.Println(s)
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	fmt.Println(t)
}
