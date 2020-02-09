package db

import (
	"fileStore/db/mysql"
	"fmt"
)

// OnFileUploadFinished returns if the file is upload to databse
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	// fmt.Printf("%s, %s, %d, %s\n", filehash, filename, filesize, fileaddr)
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file(file_sha1,file_name,file_size,file_addr, status)values(?,?,?,?,1)")

	if err != nil {
		fmt.Println("Failed to prepare statement, error: " + err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println("===>" + err.Error())
		return false
	}

	//if filesha1 already existed, <=0 means nothing changed, because filesha1 is UNIQUE KEY in tbl_file
	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before", filehash)
		}
		return true
	}

	return false
}
