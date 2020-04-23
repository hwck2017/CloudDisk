package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/Go-SQL-Driver/MySQL"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:ck@txyun123@tcp(localhost:3306)/filemeta?charset=utf8")
	db.SetMaxOpenConns(1024)
	err := db.Ping()
	if err != nil {
		fmt.Println("ping db failed, err: ", err.Error())
	}
}

// DBConn : get mysql conn
func DBConn() *sql.DB {
	return db
}
