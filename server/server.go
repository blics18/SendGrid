package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
)

// *** Structs *** 

type User struct {
	UserID int 
	Emails []string 
}

type Payload struct {
	UserData User 
}

// **************

// func serveRest(w http.ResponseWriter, r *http.Request) {
// 	response, err := getJsonResponse()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Fprintf(w, string(response))
// }

func retrieve(w http.ResponseWriter, r *http.Request) {
	//url := "http://localhost:8081"
	// res, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer res.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		panic(err)
	}

	var p Payload

	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	fmt.Println("User ID: ", p.UserData.UserID)
	fmt.Println("User Emails: ", p.UserData.Emails)

	// *** Can write back to Client here if we want

	defer r.Body.Close()
}

func main() {
	url := ":8081"
	http.HandleFunc("/retrieve", retrieve)
	log.Fatal(http.ListenAndServe(url, nil))
}
