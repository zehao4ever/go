package controller

import (
	"Server/dao"
	"encoding/json"
	//	"fmt"
	//	"log"
	//	"net"
)

const (
	RESULT_LOGIN_SUCCESS             = 1
	RESULT_LOGIN_FAIL                = 2
	RESULT_REG_SUCCESS               = 3
	RESULT_REG_FAIL                  = 4
	RESULT_MODIFY_PASS_SUCCESS       = 5
	RESULT_MODIFY_PASS_FAIL          = 6
	RESULT_MODIFY_INFO_SUCCESS       = 7
	RESULT_MODIFY_INFO_FAIL          = 8
	RESULT_MODIFY_OTHER_INFO_SUCCESS = 9
	RESULT_MODIFY_OTHER_INFO_FAIL    = 10
)

//login
func Login(arg_data TranData) (TranData, error) {
	userInfo := arg_data.UserInfo
	var user dao.User
	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		arg_data.Result = RESULT_LOGIN_FAIL
		return arg_data, err
	}
	if exist, err := dao.CheckUser(user.Uid, user.Pass); !exist || err != nil {
		arg_data.Result = RESULT_LOGIN_FAIL
		return arg_data, err
	}
	puser, err := dao.QueryUserInfo(user.Uid)
	if err != nil {
		arg_data.Result = RESULT_LOGIN_FAIL
		return arg_data, err
	}

	var buf []byte
	buf, err = json.Marshal(*puser)
	if err != nil {
		arg_data.Result = RESULT_LOGIN_FAIL
		return arg_data, err
	}
	arg_data.UserInfo = string(buf)
	arg_data.Result = RESULT_LOGIN_SUCCESS
	return arg_data, nil
}

func Register(arg_data TranData) (TranData, error) {
	userInfo := arg_data.UserInfo
	var user dao.User
	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		arg_data.Result = RESULT_LOGIN_FAIL
		return arg_data, err
	}
	id, err := dao.AddUser(user)
	if err != nil {
		arg_data.Result = RESULT_REG_FAIL
		return arg_data, err
	}
	user.Uid = id
	userBuf, err := json.Marshal(user)
	if err != nil {
		arg_data.Result = RESULT_REG_FAIL
		return arg_data, err
	}
	arg_data.SendId = id
	arg_data.UserInfo = string(userBuf)
	arg_data.Result = RESULT_REG_SUCCESS
	return arg_data, nil
}

func ModifyPass(arg_data TranData) (TranData, error) {
	userInfo := arg_data.UserInfo
	var user dao.User
	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		arg_data.Result = RESULT_MODIFY_PASS_FAIL
		return arg_data, err
	}
	if success, err := dao.CheckUser(user.Uid, user.Pass); !success || err != nil {
		arg_data.Result = RESULT_MODIFY_PASS_FAIL
		return arg_data, err
	}
	if success, err := dao.ModifyUserPass(arg_data.SendId, arg_data.Object); !success || err != nil {
		arg_data.Result = RESULT_MODIFY_PASS_FAIL
		return arg_data, err
	}
	arg_data.Result = RESULT_MODIFY_PASS_SUCCESS
	return arg_data, nil
}

func ModifyUserInfo(arg_data TranData) (TranData, error) {
	userInfo := arg_data.Object
	var user dao.User
	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		arg_data.Result = RESULT_MODIFY_INFO_FAIL
		return arg_data, err
	}

	if success, err := dao.ModifyUser(user); !success || err != nil {
		arg_data.Result = RESULT_MODIFY_INFO_FAIL
		return arg_data, err
	}
	arg_data.UserInfo = arg_data.Object
	arg_data.Result = RESULT_MODIFY_INFO_SUCCESS
	return arg_data, nil

}

func ModifyOtherUserInfo(arg_data TranData) (TranData, error) {
	return arg_data, nil
}
