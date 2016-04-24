package main

import (
	//	"encoding/json"
	"fmt"
	//"strconv"
)

type User struct {
	User_id   int    `json:"aa,string"`
	User_name string `json:"user_name,,omitempty"`
	Password  string `json:"password"`
}

func F() (int, error) {
	i := 1
	// append(i, 2)
	//	i[0] = 1
	return i, nil
}

func main() {
	// userId := 2000
	// userName := "zehao"
	//  pass := "000"
	// var u *User
	// //	userIdStr := strconv.Itoa(userId)
	// m := make(map[string]string)
	// m["aa"] = "22"
	// m["user_name"] = "zehao"
	// //	jsonStr := "{\"user_id\":\"" + userIdStr + "\",\"user_name\":\"" + userName + "\"}"
	// jsonStr, _ := json.Marshal(m)
	// if u == nil {
	// 	fmt.Println("true")
	// }
	// fmt.Println(*u)
	// json.Unmarshal(jsonStr, &u)
	// fmt.Println(u)
	// i := 0
	// i := F()
	// fmt.Println(i)
	//
	i := 0
	switch i {
	case 0:
		fmt.Println("ifgjfggbgcbgcbgjbhgkv")
		continue
		fmt.Println("fvuhhhhhhhh")
	}
}
