package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(192.168.2.23:13306)/fileserver?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, error: ", err.Error())
		os.Exit(1)
	}
}

// DBConn returns database connection object
func DBConn() *sql.DB {
	return db
}
