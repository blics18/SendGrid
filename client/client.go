package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID *int
	Email  []string
}

type HealthStatus struct {
	AppName            string
	AppVersion         string
	HealthCheckVersion string
	Port               string
	Results            struct {
		ServerStatus struct {
			OK bool
		}
		ConnectedToDB struct {
			OK bool
		}
	}
}

type HitMiss struct {
	Hits  int
	Total int
	Suppressions []string
}

func Check(userID int, emails []string) (error, HitMiss) {
	user := User{
		UserID: &userID,
		Email:  emails,
	}

	var hitMissStruct HitMiss

	userJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err, hitMissStruct
	}

	req, err := http.NewRequest("GET", "http://localhost:8082/checkBF", bytes.NewBuffer(userJSON))
	if err != nil {
		return err, hitMissStruct
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err, hitMissStruct
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, hitMissStruct
	}

	err = json.Unmarshal(body, &hitMissStruct)
	if err != nil {
		return err, hitMissStruct
	}

	fmt.Println("Response: ", string(body))

	resp.Body.Close()
	return nil, hitMissStruct
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

	var tableNames []string           // tableNames is a list of tables. Example: [User00, User01, User03, ...]
	userMap := make(map[int][]string) // userMap is a map that consists of User ID's as keys, with their value as a list of emails. //Ex: [5: ["jim@yahoo.com", "trevor@aol.com"]

	stmt := fmt.Sprintf("SELECT TABLE_NAME AS tableName FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA='UserStructs'")
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("Error from Database Connection")
	}

	// populate tableNames
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tableNames = append(tableNames, tableName)
	}

	rows.Close()

	// build userMap
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

			_, exists := userMap[id]

			if exists {
				userMap[id] = append(userMap[id], email)
			} else {
				userMap[id] = []string{email}
			}
		}

		rows.Close()
	}

	writeDataToFile(userMap)

	userList := make([]User, len(userMap)) // userList is a list of User structs: [User, User, User]
	index := 0

	// build userList from the values in userMap
	for key, value := range userMap {
		userList[index] = User{
			UserID: &key,
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
	db.Close()
	return nil
}

func HealthCheck() error {
	req, err := http.NewRequest("GET", "http://localhost:8082/healthBF", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("StatusInternalServerError: Error Code 500")
		return err
	}

	fmt.Println("Content-type: ", resp.Header["Content-Type"][0])
	fmt.Println("Date: ", resp.Header["Date"][0])
	fmt.Println("Protocol: ", resp.Proto, "\n")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	resp.Body.Close()
	return nil
}

func writeDataToFile(userMap map[int][]string) {
	file, err := os.Create("data/data.txt")
	if err != nil {
		return
	}

	defer file.Close()

	for key, value := range userMap {
		emailsNotInBF := makeRandomEmails(rand.Intn(5) + 1)
		dataString := fmt.Sprintf("%d:%s %s\n", key, strings.Join(value, " "), strings.Join(emailsNotInBF, " "))

		file.WriteString(dataString)

	}
}
