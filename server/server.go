package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/blics18/SendGrid/client"
	"github.com/cyberdelia/go-metrics-graphite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rcrowley/go-metrics"
	"github.com/willf/bloom"
)

type bloomFilter struct {
	Filter *bloom.BloomFilter
	db     *sql.DB
	cfg    client.Config
}

func NewBloomFilter(size int) *bloomFilter {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/UserStructs")
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
		db:     db,
		cfg:    cfg,
	}
}

func (bf *bloomFilter) populateBF(w http.ResponseWriter, r *http.Request) {
	timer := metrics.GetOrRegisterTimer("bloom.Filter.populateBF_Response", nil)
	start := time.Now()

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the body of the request"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the body of the request"))
		return
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

	metrics.GetOrRegisterCounter("bloom.Filter.populateBF_Request", nil).Inc(1)

	for _, user := range users {
		for _, email := range user.Email {
			bf.Filter.Add([]byte(fmt.Sprintf("%d|%s", *user.UserID, email)))
			fmt.Println(fmt.Sprintf("userID: %d", *user.UserID))
			fmt.Println(fmt.Sprintf("Email: %s", email))
		}
	}

	duration := time.Since(start)
	timer.Update(duration)
}

func (bf *bloomFilter) checkBF(w http.ResponseWriter, r *http.Request) {
	timer := metrics.GetOrRegisterTimer("bloom.Filter.checkBF_Response", nil)
	start := time.Now()

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need content to check"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the body of the request"))
		return
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

	statStruct := &client.Stats{
		UserID:			   0,
		Hits:              0,
		Miss:              0,
		NumEmails:         0,
		Suppressions:      []string{},
		TotalSuppressions: 0,
	}

	statStruct.UserID = *user.UserID

	for _, email := range user.Email {
		if bf.Filter.Test([]byte(fmt.Sprintf("%d|%s", *user.UserID, email))) {
			inDB, err := crossCheck(bf.db, bf.cfg, user.UserID, email)
			if err == nil && inDB == true {
				statStruct.Suppressions = append(statStruct.Suppressions, email)
				statStruct.TotalSuppressions += 1
				statStruct.NumEmails += 1
				statStruct.Hits += 1
			} else {
				statStruct.Miss += 1
				statStruct.NumEmails += 1
			}
		} else {
			statStruct.NumEmails += 1
			statStruct.Hits += 1
		}
	}

	statJSON, err := json.MarshalIndent(statStruct, "", " ")
	if err != nil {
		return
	}

	w.Write(statJSON)

	metrics.GetOrRegisterCounter("bloom.Filter.checkBF_Request", nil).Inc(1)

	duration := time.Since(start)
	timer.Update(duration)
}

func crossCheck(db *sql.DB, cfg client.Config, UserID *int, Email string) (bool, error) {
	var email string
	stmt := fmt.Sprintf("SELECT uid, email FROM Unsub%02d WHERE uid=? AND email=?", (*UserID)%5)
	
	err := db.QueryRow(stmt, *UserID, Email).Scan(&email)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (bf *bloomFilter) clearBF(w http.ResponseWriter, r *http.Request) {
	timer := metrics.GetOrRegisterTimer("bloom.Filter.clearBF_Response", nil)
	start := time.Now()
	
	bf.Filter.ClearAll()
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully Cleared Bloom Filter"))
	
	metrics.GetOrRegisterCounter("bloom.Filter.clearBF_Request", nil).Inc(1)

	duration := time.Since(start)
	timer.Update(duration)
}

func (bf *bloomFilter) healthBF(w http.ResponseWriter, r *http.Request) {
	timer := metrics.GetOrRegisterTimer("bloom.Filter.healthBF_Response", nil)
	start := time.Now()

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

	metrics.GetOrRegisterCounter("bloom.Filter.healthBF_Request", nil).Inc(1)

	duration := time.Since(start)
	timer.Update(duration)
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