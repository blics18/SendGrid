package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/blics18/SendGrid/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/willf/bloom"
)

type bloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter(size int) *bloomFilter {
	return &bloomFilter{
		filter: bloom.New(20*uint(size), 5),
	}
}

func (bf *bloomFilter) populateBF(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.Write([]byte("Could not read the body of the request"))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the body of the request"))
	}

	defer r.Body.Close()

	var users []client.User

	err = json.Unmarshal(body, &users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse json body"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(http.StatusOK)))

	for _, user := range users {
		for _, email := range user.Email {
			bf.filter.Add([]byte(fmt.Sprintf("%d|%s", *user.UserID, email)))
			fmt.Println(fmt.Sprintf("userID: %d", *user.UserID))
			fmt.Println(fmt.Sprintf("Email: %s", email))
		}
	}
}

func (bf *bloomFilter) checkBF(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need content to check"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the body of the request"))
	}

	defer r.Body.Close()

	var user client.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse json body"))
		return
	}

	if &user.UserID == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need User ID"))
		return
	}

	if user.Email == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need User Emails"))
		return
	}

	hitMissStruct := &client.HitMiss{
		Hits:  0,
		Total: 0,
		Suppressions: []string{},
	}

	for _, email := range user.Email {
		if bf.filter.Test([]byte(fmt.Sprintf("%d|%s", *user.UserID, email))) {
			//w.Write([]byte(email + " is in the bloom filter. Cross checking..."))
			err, inDB := crossCheck(user.UserID, email)
			if err == nil && inDB == true {
				//w.Write([]byte(email + " is in the database"))
				//fmt.Println(email + " is in the database")
				hitMissStruct.Suppressions = append(hitMissStruct.Suppressions, email) 
				hitMissStruct.Total += 1
			} else {
				//w.Write([]byte(email + " is not in the database"))
				//fmt.Println(email + " is not in the database")
				hitMissStruct.Hits += 1
				hitMissStruct.Total += 1
			}
		} else {
			//w.Write([]byte(email + " is not in the bloom filter"))
			hitMissStruct.Total += 1
		}
	}

	hitMissJSON, err := json.MarshalIndent(hitMissStruct, "", " ")
	if err != nil {
		return
	}

	w.Write(hitMissJSON)
}

func crossCheck(UserID *int, Email string) (error, bool) {
	db, err := sql.Open("mysql", "root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Println("Failed to get handle")
		db.Close()
		return err, false
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to make connection")
		db.Close()
		return err, false
	}

	var numTables int
	err = db.QueryRow("SELECT COUNT(*) AS count from information_schema.tables WHERE table_schema=?", "UserStructs").Scan(&numTables)
	if err != nil {
		fmt.Printf("Error from Database Connection")
		return err, false
	}

	stmt := fmt.Sprintf("SELECT uid, email FROM User%02d WHERE uid=? AND email=?", (*UserID)%numTables)
	rows, err := db.Query(stmt, *UserID, Email)
	
	if err != nil {
		fmt.Printf("Error from Database Connection")
		return err, false
	}

	ret := rows.Next()
	rows.Close()
	return nil, ret
}

func (bf *bloomFilter) clearBF(w http.ResponseWriter, r *http.Request) {
	bf.filter.ClearAll()
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Successfully Cleared Bloom Filter"))
}

func (bf *bloomFilter) healthBF(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	healthStruct := &client.HealthStatus{
		AppName:            "Bloom Filter",
		AppVersion:         "1.0.0",
		HealthCheckVersion: "1.0.0",
		Port:               "8082",
	}

	healthStruct.Results.ServerStatus.OK = true
	healthStruct.Results.ConnectedToDB.OK = true

	db, err := sql.Open("mysql", "root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		healthStruct.Results.ConnectedToDB.OK = false
		db.Close()
	}

	err = db.Ping()
	if err != nil {
		healthStruct.Results.ConnectedToDB.OK = false
		db.Close()
	}

	healthJSON, err := json.MarshalIndent(healthStruct, "", " ")
	if err != nil {
		return
	}

	w.Write(healthJSON)
}

func main() {
	port := ":8082"
	bf := NewBloomFilter(1000)
	http.HandleFunc("/populateBF", bf.populateBF)
	http.HandleFunc("/checkBF", bf.checkBF)
	http.HandleFunc("/clearBF", bf.clearBF)
	http.HandleFunc("/healthBF", bf.healthBF)
	log.Fatal(http.ListenAndServe(port, nil))
}
