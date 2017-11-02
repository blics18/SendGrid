/*
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Userid int      `json: userid`
	Emails []string `json: emails`
}

func add(w http.ResponseWriter, r *http.Request) {
	//	r.ParseForm()
	//	userid := r.FormValue("userid")
	//	email := r.FormValue("email")
	//	fmt.Printf("uid = %s	email=%s\n", userid, email)
	//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	//Reading Body Response
	jsn, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	//Unmarshal
	var user = []User{}
	err = json.Unmarshal(jsn, &user)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}
	fmt.Printf("Received: %v\n", user)

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

func main() {
	//	n := uint(1000)
	//	filter := bloom.New(20*n, 5) // load of 20, 5 keys
	//	filter.Add([]byte("Love"))
//
//		if filter.Test([]byte("sheila")) {
//			fmt.Printf("yes")
//		} else {
//			fmt.Printf("no")
//		}

	http.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(":8081", nil))
}*/

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
	//	fmt.Fprintf(w, "hello, %s why you no print", []byte(body))
	w.Write([]byte(body))
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	var p []User
	// Make a loop here later
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
	url := ":8081"
	http.HandleFunc("/retrieve", retrieve)
	log.Fatal(http.ListenAndServe(url, nil))
}
