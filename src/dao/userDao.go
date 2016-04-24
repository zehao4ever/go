package dao

import (
	"log"
)

/**
* 1.添加用户
* 返回用户id,若添加失败返回id = 0
 */
func AddUser(arg_user User) (int, error) {
	r := 0
	stmt, err := g_db.Prepare(INSERT_USER_SQL)
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_user.Name, arg_user.Pass, arg_user.Sex, arg_user.Age, arg_user.Intro, arg_user.Role)
	if err != nil {
		log.Println(err)
		return r, err
	}
	affected, _ := res.RowsAffected()
	var lastId int64 = 0
	if affected > 0 {
		lastId, _ = res.LastInsertId()
	}
	return int(lastId), nil
}

/**
* 2.删除用户
 */
func DelUser(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare(DEL_USER_SQL)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	affected, _ := res.RowsAffected()
	success := false
	if affected > 0 {
		success = true
	}
	return success, nil
}

/**
* 3.获取用户基本信息
* 形参：user_id
* 返回：若不存在该用户返回nil
 */
func QueryUserInfo(arg_id int) (*User, error) {
	stmt, err := g_db.Prepare("SELECT *FROM user WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	u := &User{}
	row, err := stmt.Query(arg_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&u.Uid, &u.Name, &u.Pass, &u.Sex, &u.Age, &u.Intro, &u.Role)
	}
	return u, nil
}

/**
* 4.修改用户基本信息
 */
func ModifyUser(arg_u User) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE user SET name=?,sex=?,age=?,intro=?,role_id=? WHERE user_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_u.Name, arg_u.Sex, arg_u.Age, arg_u.Intro, arg_u.Role, arg_u.Uid)
	if err != nil {
		log.Println(err)
		return false, err
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

/**
* 5.判断用户是否存在
 */
func IsUserExist(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("SELECT user_id FROM user WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	row, err := stmt.Query(arg_id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer row.Close()
	if row.Next() {
		return true, nil
	}
	return false, nil
}

/**
* 6.获取密码
 */
func QueryUserPass(arg_id int) (string, error) {
	stmt, err := g_db.Prepare("SELECT pwd FROM user WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer stmt.Close()
	pass := ""
	row := stmt.QueryRow(arg_id)
	row.Scan(&pass)
	return pass, nil
}

/**
* 7.修改密码
 */
func ModifyUserPass(arg_id int, arg_pass string) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE user SET pwd = ? WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_pass, arg_id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

/**
* 8.修改用户角色
* 形参：user_id,role_id
 */
func ModifyUserRole(arg_id int, arg_role int) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE user SET role_id = ? WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_role, arg_id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

/**
* 9.获取用户角色类型
 */
func QueryUserRole(arg_id int) (int, error) {
	stmt, err := g_db.Prepare("SELECT role_id FROM user WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()
	role := 0
	row := stmt.QueryRow(arg_id)
	row.Scan(&role)
	return role, nil
}

/**
* 10.获取用户权限
 */
func QueryUserAuth(arg_role int) ([]int, error) {
	auths := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT auth_id FROM role_auth WHERE role_id = ?")
	if err != nil {
		log.Println(err)
		return auths, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(arg_role)
	if err != nil {
		log.Println(err)
		return auths, err
	}
	defer rows.Close()
	v := 0
	for rows.Next() {
		rows.Scan(&v)
		auths = append(auths, v)
	}
	return auths, nil
}
