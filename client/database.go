package client

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func insertToTables(numTables int, usr User, db *sql.DB) error {
	stmt := fmt.Sprintf("INSERT INTO Unsub%02d(uid,email) VALUE(?, ?)", *(usr.UserID)%numTables)
	stmtHandle, err := db.Prepare(stmt)
	if err != nil {
		return err
	}

	defer stmtHandle.Close()

	for _, email := range usr.Email{
		_, err := stmtHandle.Exec(usr.UserID, email)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func PopulateDB(numUsers int, numEmails int, numTables int) (*sql.DB, error) {
	p := MakeRandomUsers(numUsers, numEmails)

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
		return nil, err
	}

	for i := 0; i < numUsers; i++ {
		err := insertToTables(numTables, p[i], db)
		if err != nil {
			fmt.Println(err)
			db.Close()
		}
	}

	return db, nil
}

func DropTables(numTables int, db *sql.DB) error {
	for i := 0; i < numTables; i++ {
		stmt := fmt.Sprintf("DROP TABLE Unsub%02d", i)
		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}