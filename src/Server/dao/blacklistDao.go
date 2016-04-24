package dao

import (
	"log"
)

/**
* 1.获取所有黑名单用户
* 返回：元素为user_id的整型数组
 */
func QueryBlacklist() ([]int, error) {
	ids := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT user_id FROM blacklist")
	if err != nil {
		log.Println(err)
		return ids, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
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
* 2.把用户加入黑名单
* 形参：user_id
 */
func InsertIntoBlacklist(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("INSERT INTO blacklist (user_id,join_time) VALUES(?,null)")
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
	if affected > 0 {
		return true, nil
	}
	return false, nil
}

/**
* 3.把用户移除黑名单
* 形参：user_id
 */
func DelFromBlacklist(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM blacklist WHERE user_id = ?")
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
* 3.判断指定用户是否在黑名单中
 */
func IsInBlacklist(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("SELECT user_id FROM blacklist WHERE user_id = ?")
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
* 4.更新指定用户在黑名单中的时间
* 形参：user_id
 */
func UpdateTimeInBlacklist(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("UPDATE blacklist SET join_time = null WHERE user_id=?")
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
	if affected > 0 {
		return true, nil
	}
	return false, nil
}
