package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var instanceDB *sql.DB

func ConnectDatabase() *sql.DB {
	//dbUser := os.Getenv("DB_USERNAME")
	//dbPassword := os.Getenv("DB_PASSWORD")
	//dbDatabase := os.Getenv("DATABASE_NAME")
	//dbHost := os.Getenv("DATABASE_HOST")
	//
	//urlConnect := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", dbUser, dbPassword, dbHost, dbDatabase)

	urlConnect := "root:root@tcp(localhost:3306)/db_marketplace"

	if instanceDB != nil {
		if err := instanceDB.Ping(); err != nil {
			instanceDB = nil
		}
	}

	if instanceDB == nil {
		db, err := sql.Open("mysql", urlConnect)

		if err != nil {
			panic(err.Error())
		}
		instanceDB = db
	}

	return instanceDB
}
