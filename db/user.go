package db

import (
	"fileStore/db/mysql"
	"fmt"
)

// UserSignUp create a user when user sign up
func UserSignUp(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}

	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false
}
