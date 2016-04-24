package dao

import (
	"log"
)

/**
* 1.添加公告
* 形参：Billboard中的两个时间中的create_time为空
 */
func InsertBillboard(arg_bill Billboard) (int, error) {
	r := 0
	stmt, err := g_db.Prepare("INSERT INTO billboard (content,user_id,user_name) VALUES (?,?,?)")
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(arg_bill.Content, arg_bill.User_id, arg_bill.User_name)
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
* 2.删除公告
 */
func DelBillboard(arg_id int) (bool, error) {
	stmt, err := g_db.Prepare("DELETE FROM billboard WHERE bill_id = ?")
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
* 4.查询公告
 */
func QueryBillboard(arg_id int) (*Billboard, error) {
	stmt, err := g_db.Prepare("SELECT *FROM billboard WHERE bill_id = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	b := &Billboard{}
	row, err := stmt.Query(arg_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&b.Bid, &b.Content, &b.CreateTime, &b.User_id, &b.User_name)
	}
	return b, nil
}

/**
* 获取用户所有的公告id
 */
func QueryBillboardIdOfUser(arg_id int) ([]int, error) {
	ids := make([]int, 0)
	stmt, err := g_db.Prepare("SELECT bill_id FROM billboard WHERE user_id = ?")
	if err != nil {
		log.Println(err)
		return ids, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(arg_id)
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
