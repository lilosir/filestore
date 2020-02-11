package db

import (
	"fileStore/db/mydb"
	"fileStore/db/mysql"
	"fmt"
)

// UserSignUp create a user when user sign up
func UserSignUp(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false
}

// UserSignIn check user if login successfully
func UserSignIn(username, encpassword string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpassword {
		return true
	}

	return false
}

// UpdateToken will create/update the token from tbl_use_token table
func UpdateToken(username, token string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`, `user_token`) values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
