package dao

import (
	"log"
)

/**
* 1.添加好友
* 形参：user_id,user_id
 */
func InsertFriend(arg_f int, arg_s int) (bool, error) {
	stmt, err := g_db.Prepare("INSERT INTO friendship VALUES(?,?,null)")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_f, arg_s)
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
* 2.删除好友
 */
func DelFriend(arg_f int, arg_s int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM friendship WHERE first_user_id = ?,second_id=?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_f, arg_s)
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
* 3.查询所有好友
 */
func QueryFriends(arg_f int) ([]int, error) {
	ids := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT second_user_id FROM friendship WHERE first_user_id=?")
	if err != nil {
		log.Println(err)
		return ids, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(arg_f)
	if err != nil {
		log.Println(err)
		return ids, err
	}
	defer rows.Close()
	v := 0
	for rows.Next() {
		rows.Scan(&v)
		ids = append(ids, v)
	}
	return ids, nil
}

/**
* 4.获取建立好友关系的时间
* 返回的时间为字符串：xxxx-xx-xx xx:xx:xx
 */
func QueryBeFriendTime(arg_f int, arg_s int) (string, error) {
	stmt, err := g_db.Prepare("SELECT create_time FROM friendship WHERE first_user_id = ? AND second_user_id =?")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer stmt.Close()
	row, err := stmt.Query(arg_f, arg_s)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer row.Close()
	t := ""
	for row.Next() {
		row.Scan(&t)
	}
	return t, nil
}
