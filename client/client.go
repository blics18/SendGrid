package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/caarlos0/env"
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
	Hits         int
	Miss         int
	Total        int
	Suppressions int
}

type Config struct {
	Size             int    `env:"BLOOM_SIZE" envDefault:"1000"`
	Port             string `env:"BLOOM_PORT" envDefault:"8082"`

	NumTables        int    `env:"BLOOM_NUM_TABLES" envDefault:"5"`
	NumUsers         int    `env:"BLOOM_NUM_USERS" envDefault:"10"`
	NumEmails        int    `env:"BLOOM_NUM_EMAILS" envDefault:"1000"`
	NumHashFunctions uint   `env:"BLOOM_NUM_HASH_FUNCTIONS envDefault:"5"`
}

func Check(cfg Config, userID int, emails []string) (error, HitMiss) {
	user := User{
		UserID: &userID,
		Email:  emails,
	}

	var hitMissStruct HitMiss

	userJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err, hitMissStruct
	}

	endpoint := fmt.Sprintf("http://localhost:%s/checkBF", cfg.Port)

	req, err := http.NewRequest("GET", endpoint, bytes.NewBuffer(userJSON))
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

func Clear(cfg Config) error {
	endpoint := fmt.Sprintf("http://localhost:%s/clearBF", cfg.Port)
	req, err := http.NewRequest("GET", endpoint, nil)
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

func Populate(cfg Config) error {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		db.Close()
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
		return err
	}

	var tableNames []string           // tableNames is a list of tables. Example: [User00, User01, User03, ...]
	userMap := make(map[int][]string) // userMap is a map that consists of User ID's as keys, with their value as a list of emails. //Ex: [5: ["jim@yahoo.com", "trevor@aol.com"]

	stmt := fmt.Sprintf("SELECT TABLE_NAME AS tableName FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_SCHEMA='UserStructs'")
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("Error from Database Connection")
		return err
	}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tableNames = append(tableNames, tableName)
	}

	rows.Close()

	for _, tableName := range tableNames {
		stmt := fmt.Sprintf("SELECT uid, email FROM UserStructs.%s", tableName)
		rows, err := db.Query(stmt)
		if err != nil {
			fmt.Printf("Error from Database Connection")
			return err
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

	userList := make([]User, len(userMap))
	index := 0

	for key, value := range userMap {
		temp := key
		userList[index] = User{
			UserID: &temp,
			Email:  value,
		}
		index += 1
	}

	userJSON, err := json.MarshalIndent(userList, "", "  ")
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("http://localhost:%s/populateBF", cfg.Port)

	req, err := http.NewRequest("GET", endpoint, bytes.NewBuffer(userJSON))
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

func HealthCheck(cfg Config) error {
	endpoint := fmt.Sprintf("http://localhost:%s/healthBF", cfg.Port)

	req, err := http.NewRequest("GET", endpoint, nil)
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
		for _, email := range value {
			dataString := fmt.Sprintf("%d:%s\n", key, email)
			file.WriteString(dataString)
		}
		emailsNotInBF := makeRandomEmails(rand.Intn(1000) + 1)
		for _, email := range emailsNotInBF {
			dataString := fmt.Sprintf("%d:%s\n", key, email)
			file.WriteString(dataString)
		}
	}
}

func GetEnv() Config {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return cfg
}
