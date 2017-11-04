package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/willf/bloom"
)

// *** GLOBAL VARIABLES ***

var bloomFilter *bloom.BloomFilter

// *** STRUCTS ***

type User struct {
	UserID int
	Email  []string
}

func createBloomFilter() *bloom.BloomFilter {
	n := uint(1000)
	filter := bloom.New(20*n, 5)
	return filter
}

func populateBF(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	var p []User

	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse json body"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))

	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i].Email); j++ {
			bloomFilter.Add([]byte(fmt.Sprintf("%d|%s", p[i].UserID, p[i].Email[j])))
			fmt.Println(fmt.Sprintf("userID: %d", p[i].UserID))
			fmt.Println(fmt.Sprintf("Email: %s", p[i].Email[j]))
		}
	}

}

func checkBF(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need content to check"))
		return

	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	var p User

	err = json.Unmarshal(body, &p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse json body"))
		return
	}
	fmt.Println(p.UserID)
	if &p.UserID == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need User ID"))
		return
	}

	if p.Email == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Need User Emails"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))

	for i := 0; i < len(p.Email); i++ {

		if bloomFilter.Test([]byte(fmt.Sprintf("%d|%s", p.UserID, p.Email[i]))) {
			fmt.Print(p.Email[i])
			fmt.Println(" is in the bloom filter.")
		} else {
			fmt.Println("NO")
		}
	}

}

func clearBF(w http.ResponseWriter, r *http.Request) {
	bloomFilter.ClearAll()

}

func main() {
	url := ":8082"
	bloomFilter = createBloomFilter()
	http.HandleFunc("/populateBF", populateBF)
	http.HandleFunc("/checkBF", checkBF)
	http.HandleFunc("/clearBF", clearBF)
	log.Fatal(http.ListenAndServe(url, nil))
}
