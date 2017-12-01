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
	Filter *bloom.BloomFilter
	db *sql.DB
	cfg client.Config
}

func NewBloomFilter(size int) *bloomFilter {
	db, err := sql.Open("mysql", "root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Println("Failed to get handle")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to make connection")
	}

	cfg := client.GetEnv()

	return &bloomFilter{
		Filter: bloom.New(20*uint(cfg.Size), cfg.NumHashFunctions),
		db: db,
		cfg: cfg,
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

	for _, user := range users {
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
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(http.StatusOK)))

	for _, user := range users {
		for _, email := range user.Email {
			bf.Filter.Add([]byte(fmt.Sprintf("%d|%s", *user.UserID, email)))
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
		Hits:  	      0,
		Miss:         0,
		Total:        0,
		Suppressions: 0,
	}

	// var suppresions []string

	for _, email := range user.Email {
		if bf.Filter.Test([]byte(fmt.Sprintf("%d|%s", *user.UserID, email))) {
			//w.Write([]byte(email + " is in the bloom filter. Cross checking..."))
			err, inDB := crossCheck(bf.db, bf.cfg, user.UserID, email)
			if err == nil && inDB == true {
				//w.Write([]byte(email + " is in the database"))
				//fmt.Println(email + " is in the database")
				//suppresions = append(suppresions, email)
				hitMissStruct.Suppressions += 1
				hitMissStruct.Total += 1
				hitMissStruct.Hits += 1
			} else {
				//w.Write([]byte(email + " is not in the database"))
				//fmt.Println(email + " is not in the database")
				hitMissStruct.Miss += 1
				hitMissStruct.Total += 1
			}
		} else {
			//w.Write([]byte(email + " is not in the bloom filter"))
			//fmt.Println(email + " is not in the BF")
			hitMissStruct.Total += 1
			hitMissStruct.Hits += 1
		}
	}

	// hitMissStruct.Suppressions = len(suppresions)

	hitMissJSON, err := json.MarshalIndent(hitMissStruct, "", " ")
	if err != nil {
		return
	}

	w.Write(hitMissJSON)
}

func crossCheck(db *sql.DB, cfg client.Config, UserID *int, Email string) (error, bool) {
	stmt := fmt.Sprintf("SELECT uid, email FROM User%02d WHERE uid=? AND email=?", (*UserID)%cfg.NumTables)
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
	bf.Filter.ClearAll()
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Successfully Cleared Bloom Filter"))
}

func (bf *bloomFilter) healthBF(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	healthStruct := &client.HealthStatus{
		AppName:            "Bloom Filter",
		AppVersion:         "1.0.0",
		HealthCheckVersion: "1.0.0",
		Port:               bf.cfg.Port,
	}

	healthStruct.Results.ServerStatus.OK = true
	healthStruct.Results.ConnectedToDB.OK = true

	err := bf.db.Ping()
	if err != nil {
		healthStruct.Results.ConnectedToDB.OK = false
	}

	healthJSON, err := json.MarshalIndent(healthStruct, "", " ")
	if err != nil {
		return
	}

	w.Write(healthJSON)
}

func main() {
	cfg := client.GetEnv()
	bf := NewBloomFilter(cfg.Size)
	http.HandleFunc("/populateBF", bf.populateBF)
	http.HandleFunc("/checkBF", bf.checkBF)
	http.HandleFunc("/clearBF", bf.clearBF)
	http.HandleFunc("/healthBF", bf.healthBF)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
