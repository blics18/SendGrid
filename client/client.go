package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID *int
	Email  []string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Check(userID int, emails []string) error {
	user := User{
		UserID: &userID,
		Email:  emails,
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", "http://localhost:8082/checkBF", bytes.NewBuffer(userJSON))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}

func Clear() error {
	req, err := http.NewRequest("GET", "http://localhost:8082/clearBF", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}

func Populate() error {
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

	var tableNames []string
	userMap := make(map[*int][]string)

	stmt := fmt.Sprintf("SELECT TABLE_NAME AS tableName FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA='UserStructs'")
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Printf("Error from Database Connection")
	}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tableNames = append(tableNames, tableName)
	}

	for _, tableName := range tableNames {
		stmt := fmt.Sprintf("SELECT uid, email FROM UserStructs.%s", tableName)
		rows, err := db.Query(stmt)
		if err != nil {
			fmt.Printf("Error from Database Connection")
		}

		for rows.Next() {
			var id int
			var email string
			
			rows.Scan(&id, &email)
			_, exists := userMap[&id]
			if exists {
				userMap[&id] = append(userMap[&id], email)
			} else {
				userMap[&id] = []string{email}
			}
		}
	}

	userList := make([]User, len(userMap))
	index := 0

	for key, value := range userMap {
		userList[index] = User{
			UserID: key,
			Email:  value,
		}
		index++
	}

	userJSON, err := json.MarshalIndent(userList, "", "  ")

	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "http://localhost:8082/populateBF", bytes.NewBuffer(userJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}