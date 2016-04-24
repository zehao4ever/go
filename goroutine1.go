package main

import (
	"fmt"
	// "time"
	// "runtime"
	"encoding/json"
	//	"sync"
)

type User struct {
	Id int   `json:"id,string"`
	F  []int `json:"list"`
}

func main() {
	var u User
	u.Id = 110
	u.F = make([]int, 2)
	u.F[0] = 1
	u.F[1] = 2
	buf, _ := json.Marshal(u)
	fmt.Println(string(buf))
	var v User
	json.Unmarshal(buf, &v)
	fmt.Println(v)
}
