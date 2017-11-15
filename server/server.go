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

func createBloomFilter() *bloomFilter {
	return &bloomFilter{
		filter: bloom.New(20*uint(1000), 5),
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

	for _, email := range user.Email {
		if bf.filter.Test([]byte(fmt.Sprintf("%d|%s", *user.UserID, email))) {
			w.Write([]byte(email + " is in the bloom filter. Cross checking..."))
			if crossCheck(user.UserID, email) {
				w.Write([]byte(email + " is in the database"))
				fmt.Println(email + " is in the database")
			} else {
				w.Write([]byte(email + " is not in the database"))
				fmt.Println(email + " is not in the database")
			}

		} else {
			w.Write([]byte(email + " is not in the bloom filter"))
		}
	}
}

func crossCheck(UserID *int, Email string) bool {
	db, err := sql.Open("mysql", "root:SendGrid@tcp(localhost:3306)/UserStructs")
	if err != nil {
		fmt.Printf("Failed to get handle\n")
		db.Close()
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Unable to make connection\n")
		db.Close()
	}
	const numTables int = 5
	stmt := fmt.Sprintf("SELECT uid, email FROM User%02d WHERE uid=%d AND email='%s'", (*UserID)%numTables, *UserID, Email)
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Printf("Error from Database Connection")
		return false
	}
	return rows.Next()

}
func (bf *bloomFilter) clearBF(w http.ResponseWriter, r *http.Request) {
	bf.filter.ClearAll()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully Cleared Bloom Filter"))
}

func main() {
	port := ":8082"
	bf := createBloomFilter()
	http.HandleFunc("/populateBF", bf.populateBF)
	http.HandleFunc("/checkBF", bf.checkBF)
	http.HandleFunc("/clearBF", bf.clearBF)
	log.Fatal(http.ListenAndServe(port, nil))
}
