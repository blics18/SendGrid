package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// *** Structs ***

type User struct {
	UserID int
	Email  []string
}

type Payload struct {
	UserData []User
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
	
	defer r.Body.Close()


	var p []User

	err = json.Unmarshal(body, &p)
	if err != nil {
	    w.WriteHeader(http.StatusBadRequest)
	    w.Write([]byte("Unable to parse json body"))
	    return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("success"))

	// Make a loop here later
	for i := 0; i < len(p); i++ {
		fmt.Println("User ID: ", p[i].UserID)
		fmt.Println("User Email: ", p[i].Email)
	}

	// *** Can write back to Client here if we want
}

func main() {
	//I have mine on 8082 because my 8081 Server wouldn't shut down
	url := ":8082"
	http.HandleFunc("/retrieve", retrieve)
	log.Fatal(http.ListenAndServe(url, nil))
}
