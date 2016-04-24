package controller

import "net"

type TranData struct {
	ActionType int      `json:"actionType"`
	UserInfo   string   `json:"userInfo"`
	Result     int      `json:"result"`
	SendId     int      `json:"sendId"`
	RecvId     int      `json:"recvId"`
	Object     string   `json:"object"`
	Addr       string   `json:"-"`
	Conn       net.Conn `json:"-"`
}
