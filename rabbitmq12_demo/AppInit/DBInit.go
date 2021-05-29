package AppInit

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func DBInit() error {
	var err error
	db, err := sqlx.Connect("mysql",
		"root:root@tcp(localhost:3306)/Order?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		return err
	}
	db.SetMaxIdleConns(20)
	db.SetMaxIdleConns(10)
	return nil
}

func GetDB() *sqlx.DB {
	return db
}
