package models

import (
	"database/sql"
	"fmt"
)

// AdminUser struct
type AdminUser struct {
	Id       int    `db:"admin_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

//查询
func QueryRowDB(sql string) *sql.Row {
	return Db.QueryRow(sql)
}

//根据用户名查询id
func QueryUserWithUsername(username string) int {
	sql := fmt.Sprintf("where username='%s'", username)
	return QueryUserWightCon(sql)
}

//根据用户名和密码，查询id
func QueryUserWithParam(username, password string) int {
	sql := fmt.Sprintf("where username='%s' and password='%s'", username, password)
	return QueryUserWightCon(sql)
}

//按条件查询
func QueryUserWightCon(con string) int {
	sql := fmt.Sprintf("select admin_id from tbl_admin %s", con)
	fmt.Println(sql)
	row := QueryRowDB(sql)
	id := 0
	row.Scan(&id)
	return id
}
