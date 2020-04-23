package db

import (
	mydb "CloudDisk/db/mysql"
	"fmt"
)

//InsertIntoDB insert filemeta to db
func InsertIntoDB(fileSha1, filename, filePath string, fileSize int64) bool {
	db := mydb.DBConn()
	stmt, _ := db.Prepare("insert into tbl_file(`file_sha1`, `file_name`, `file_size`, " +
		"`file_path`, `status`) values (?,?,?,?,1)")

	defer stmt.Close()
	res, err := stmt.Exec(fileSha1, filename, fileSize, filePath)
	if err != nil {
		fmt.Println("insert metadata to sql failed")
		return false
	}

	n, err := res.RowsAffected()
	if err == nil {
		if n <= 0 {
			fmt.Printf("file %s has been insertted into mysql", filename)
		}

		return true
	}

	return false
}

//QueryFromDB query metadata from db by fileSha1
func QueryFromDB(fileSha1 string) bool {
	db := mydb.DBConn()
	stmt, _ := db.Prepare(`SELECT * From tbl_file where file_sha1 = ?`)
	defer stmt.Close()
	_, err := stmt.Exec(fileSha1)
	if err != nil {
		fmt.Printf("query metadata by %s failed", fileSha1)
		return false
	}

	return true
}
