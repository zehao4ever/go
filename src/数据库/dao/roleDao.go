package dao

import (
	"log"
)

/**
* 1.新增角色
 */
func InsertRole(arg_role Role) (int, error) {
	r := 0
	stmt, err := g_db.Prepare("INSERT INTO role (role_name,role_desc) VALUES(?,?)")
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_role.Name, arg_role.Desc)
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
* 2.删除角色
 */
func DelRole(arg_rid int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM role WHERE role_id = ?")
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
	delFlag, err := DelAuthsOfRole(arg_rid)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	if affected > 0 && delFlag {
		return true, nil
	}
	return false, nil
}

/**
* 3.获取角色对象
 */
func QueryRole(arg_rid int) (*Role, error) {
	stmt, err := g_db.Prepare("SELECT *FROM role WHERE role_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(arg_rid)
	r := &Role{}
	err = row.Scan(&r.Rid, &r.Name, &r.Desc)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r.Auths, err = QueryRoleAuths(arg_rid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return r, nil
}

/**
* 5.修改角色权限
* 形参：role_id,元素为权限id的整型数组
 */
func ModifyRoleAuth(arg_rid int, arg_aids []int) (bool, error) {
	delFlag, err := DelAuthsOfRole(arg_rid)
	if err != nil {
		log.Println(err)
		return false, err
	}
	//为权限添加要修改的角色
	addFlag, err := AddAuthsForRole(arg_rid, arg_aids)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if delFlag && addFlag {
		return true, nil
	}
	return false, nil
}

/**
* 6.获取角色权限
* 形参：role_id
 */
func QueryRoleAuths(arg_rid int) ([]int, error) {
	auths := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT auth_id FROM role_auth WHERE role_id = ?")
	if err != nil {
		log.Println(err)
		return auths, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(arg_rid)
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

/**
* 7.修改角色描述，简介
 */
func ModifyRoleInfo(arg_role Role) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE role SET role_name=?,role_desc=? WHERE role_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_role.Name, arg_role.Desc, arg_role.Rid)
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
