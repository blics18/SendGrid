package main

import (
	"SendGrid/utils"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const numTables int = 5

func main() {
	numEmails := 100
	numUsers := 10
	p := utils.MakeRandomUsers(numUsers, numEmails)
	db, err := sql.Open("mysql",
		"jose:gomez@tcp(127.0.0.1:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		db.Close()
	}
	defer db.Close()

	//Validate DSN data
	err = db.Ping()
	if err != nil {
		fmt.Printf("Unable to make connection\n")
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
	/*	//Print out our table
		rows, err := db.Query("SELECT uid, email FROM Users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&uid, &email)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(uid, email)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	*/
}

//Should I make it return anything?
func create_table(numTable int, db *sql.DB) error {
	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS User%02d (
		id int NOT NULL AUTO_INCREMENT,
		uid int(255) NOT NULL,
		email varchar(255) DEFAULT NULL, PRIMARY KEY (id)) 
		DEFAULT CHARACTER SET utf8`, numTable)
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("Failed to execute table")
		return err
	}
	/*stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS User (uid int(255) NOT NULL,email varchar(255) DEFAULT NULL, PRIMARY KEY (uid)) DEFAULT CHARACTER SET utf8")
	if err != nil {
		fmt.Println("Failed to prepare")
	}

		lastId, err := res.LastInsertId()
		fmt.Println(lastId)
	*/
	return nil
}

func insert(usr utils.User, db *sql.DB) error {
	stmt := fmt.Sprintf("INSERT INTO User%02d(uid,email) VALUE(?, ?)", usr.UserID%numTables)
	stmtHandle, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	for i := 0; i < len(usr.Email); i++ {
		_, err := stmtHandle.Exec(usr.UserID, usr.Email[i])
		if err != nil {
			fmt.Println("Failed to execute insert")
			return err
		}
	}
	return nil
}
