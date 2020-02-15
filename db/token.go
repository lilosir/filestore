package db

import (
	"fileStore/db/mysql"
	"fmt"
)

// GetUserToken return user token from databse
func GetUserToken(username string) (string, error) {
	token := ""
	stmt, err := mysql.DBConn().Prepare("select user_token from tbl_user_token where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return token, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&token)
	if err != nil {
		fmt.Println(err.Error())
		return token, err
	}
	return token, nil
}
