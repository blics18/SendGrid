package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Userid string   `json: userid`
	Emails []string `json: emails`
}

func add(w http.ResponseWriter, r *http.Request) {
	/*	r.ParseForm()
		userid := r.FormValue("userid")
		email := r.FormValue("email")
		fmt.Printf("uid = %s	email=%s\n", userid, email)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	*/

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
	log.Printf("Received: %v\n", user)
	log.Print(user[0].Userid)

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

func main() {
	/*	n := uint(1000)
		filter := bloom.New(20*n, 5) // load of 20, 5 keys
		filter.Add([]byte("Love"))

		if filter.Test([]byte("sheila")) {
			fmt.Printf("yes")
		} else {
			fmt.Printf("no")
		}

		i := uint32(100)
		n1 := make([]byte, 4)
		binary.BigEndian.PutUint32(n1, i)
		filter.Add(n1)
		if filter.EstimateFalsePositiveRate(1000) > 0.001 {
			fmt.Printf("yes")
		} else {
			fmt.Printf("no")
		}
	*/
	http.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
