package db

import (
	"database/sql"
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

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta return filemeta from mysql
func GetFileMeta(filehash string) (*TableFile, error) {

	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1=? and status=1 limit 1")

	if err != nil {
		fmt.Println("Failed to prepare statement, error: " + err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Println("Failed to query statement, error: " + err.Error())
		return nil, err
	}

	return &tfile, nil
}
