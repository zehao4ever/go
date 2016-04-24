package dao

import (
	"log"
)

/**
* 1.新增权限
 */
func InsertAuth(arg_auth Auth) (int, error) {
	r := 0
	stmt, err := g_db.Prepare("INSERT INTO auth (auth_name,auth_desc) VALUES(?,?)")
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_auth.Name, arg_auth.Desc)
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
* 2.删除权限
* 形参：auth_id
 */
func DelAuth(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM auth WHERE auth_id = ?")
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
* 3.获取拥有指定权限的全部角色
 */
func QueryAuthRoles(arg_auth int) ([]int, error) {
	roles := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT role_id FROM role_auth WHERE auth_id = ?")
	if err != nil {
		log.Println(err)
		return roles, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(arg_auth)
	if err != nil {
		log.Println(err)
		return roles, err
	}
	defer rows.Close()
	v := 0
	for rows.Next() {
		rows.Scan(&v)
		roles = append(roles, v)
	}
	return roles, nil
}

/**
* 4.修改权限
 */
func ModifyAuth(arg_auth Auth) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE auth SET auth_name=?,auth_desc=? WHERE auth_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_auth)
	if err != nil {
		log.Println(err)
		return false, err
	}
	affected, _ := res.RowsAffected()
	//先删除指定角色——权限所有记录
	delFlag, err := DelRolesOfAuth(arg_auth.Aid)
	if err != nil {
		log.Println(err)
		return false, err
	}
	//为权限添加要修改的角色
	addFlag, err := AddRolesForAuth(arg_auth.Aid, arg_auth.Role_id)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if affected > 0 && delFlag && addFlag {
		return true, nil
	}
	return false, nil
}

/**
* 5.获取权限对象
 */
func QueryAuth(arg_aid int) (*Auth, error) {
	stmt, err := g_db.Prepare("SELECT *FROM auth WHERE auth_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(arg_aid)
	auth := &Auth{}
	err = row.Scan(&auth.Aid, &auth.Name, &auth.Desc)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	auth.Role_id, err = QueryAuthRoles(arg_aid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return auth, nil
}

//删除角色对应的指定权限
func DelAuthOfRole(arg_rid, arg_aid int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM role_auth WHERE role_id=? AND auth_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_rid, arg_aid)
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

//角色增加权限
func AddAuthsForRole(arg_rid int, arg_aids []int) (bool, error) {
	stmt, err := g_db.Prepare("INSERT INTO role_auth (role_id,auth_id) VALUES(?,?)")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	sum := 0
	for i := 0; i < len(arg_aids); i++ {
		res, err := stmt.Exec(arg_rid, arg_aids[i])
		if err != nil {
			log.Println(err)
			return false, err
		}
		affected, _ := res.RowsAffected()
		sum += int(affected)
	}

	if sum == len(arg_aids) {
		return true, nil
	}
	return false, nil
}

//删除指定权限与所有角色的对应记录
func DelRolesOfAuth(arg_aid int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM role_auth WHERE auth_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_aid)
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

//删除指定角色所有的权限之间对应记录
func DelAuthsOfRole(arg_rid int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM role_auth WHERE role_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_rid)
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

//把指定权限分配给指定角色
func AddRolesForAuth(arg_aid int, arg_rids []int) (bool, error) {
	stmt, err := g_db.Prepare("INSERT INTO role_auth (role_id,auth_id) VALUES(?,?)")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	sum := 0
	for i := 0; i < len(arg_rids); i++ {
		res, err := stmt.Exec(arg_rids[i], arg_aid)
		if err != nil {
			log.Println(err)
			return false, err
		}
		affected, _ := res.RowsAffected()
		sum += int(affected)
	}
	if sum == len(arg_rids) {
		return true, nil
	}
	return false, nil
}
