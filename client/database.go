package client

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(numTables int, db *sql.DB) error {
	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS User%02d (
		id int NOT NULL AUTO_INCREMENT,
		uid int(255) NOT NULL,
		email varchar(255) DEFAULT NULL, PRIMARY KEY (id)) 
		DEFAULT CHARACTER SET utf8`, numTables)
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func insertToTables(numTables int, usr User, db *sql.DB) error {
	stmt := fmt.Sprintf("INSERT INTO User%02d(uid,email) VALUE(?, ?)", *(usr.UserID)%numTables)
	stmtHandle, err := db.Prepare(stmt)
	if err != nil {
		return err
	}

	defer stmtHandle.Close()

	for i := 0; i < len(usr.Email); i++ {
		_, err := stmtHandle.Exec(usr.UserID, usr.Email[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func PopulateDB(numUsers int, numEmails int, numTables int) *sql.DB {
	p := MakeRandomUsers(numUsers, numEmails)

	db, err := sql.Open("mysql", "root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		db.Close()
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
	}

	for i := 0; i < numTables; i++ {
		err := createTable(i, db)
		if err != nil {
			fmt.Println(err)
			db.Close()
		}
	}

	for i := 0; i < numUsers; i++ {
		err := insertToTables(numTables, p[i], db)
		if err != nil {
			fmt.Println(err)
			db.Close()
		}
	}

	return db
}

func DropTables(numTables int, db *sql.DB) error {
	for i := 0; i < numTables; i++ {
		stmt := fmt.Sprintf("DROP TABLE User%02d", i)
		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}