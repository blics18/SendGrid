package client

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const numTables int = 5

func PopulateDB() {
	numEmails := 100
	numUsers := 10
	p := MakeRandomUsers(numUsers, numEmails)
	db, err := sql.Open("mysql",
		"root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		db.Close()
	}
	defer db.Close()

	//Validate DSN data
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
	}

	for i := 0; i < numTables; i++ {
		err := create_table(i, db)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//Create Tables and Insert Data
	for i := 0; i < numUsers; i++ {
		insert(p[i], db)
	}
}

func create_table(numTable int, db *sql.DB) error {
	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS User%02d (
		id int NOT NULL AUTO_INCREMENT,
		uid int(255) NOT NULL,
		email varchar(255) DEFAULT NULL, PRIMARY KEY (id)) 
		DEFAULT CHARACTER SET utf8`, numTable)
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func insert(usr User, db *sql.DB) error {
	stmt := fmt.Sprintf("INSERT INTO User%02d(uid,email) VALUE(?, ?)", *(usr.UserID)%numTables) //take out *in usr.UserID jose
	stmtHandle, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	for i := 0; i < len(usr.Email); i++ {
		_, err := stmtHandle.Exec(usr.UserID, usr.Email[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
