package main

import (
	"dao"
	"fmt"
	"log"
)

func TestInsertUser() {
	user1 := dao.User{
		Uid:   1,
		Name:  "u1",
		Pass:  "11",
		Sex:   0,
		Age:   18,
		Intro: "i am u1",
		Role:  0,
	}
	user2 := dao.User{
		Uid:   1,
		Name:  "u2",
		Pass:  "22",
		Sex:   1,
		Age:   10,
		Intro: "i am u2",
		Role:  0,
	}
	id1, _ := dao.AddUser(user1)
	id2, _ := dao.AddUser(user2)
	fmt.Println(id1)
	fmt.Println(id2)
}

//好友关系测试
func TestAddFriend() {
	f, err := dao.InsertFriend(1, 2)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("addFriend: ", f)
	f, err = dao.InsertFriend(1, 3)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("addFriend: ", f)
	f, err = dao.InsertFriend(2, 1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("addFriend: ", f)
	f, err = dao.InsertFriend(3, 1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("addFriend: ", f)
	ids, err := dao.QueryFriends(1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("1's Friend: ", ids)
}
func main() {
	dao.InitDB()
	TestAddFriend()
	//	TestInsertUser()
	// user2 := dao.User{
	// 	Uid:   1,
	// 	Name:  "u2",
	// 	Pass:  "22",
	// 	Sex:   1,
	// 	Age:   10,
	// 	Intro: "i am u2",
	// 	Role:  0,
	// }
	// dao.ModifyUser(user2)
	// u, _ := dao.QueryUserInfo(1)
	// fmt.Println(*u)
	// s, _ := dao.QueryUserPass(1)
	// fmt.Println(s)
	// f, _ := dao.IsUserExist(1)
	// fmt.Println(f)
	// dao.DelUser(2)
	// dao.UpdateTimeInBlacklist(1)
	// f, _ := dao.IsInBlacklist(1)
	// fmt.Println(f)
}
