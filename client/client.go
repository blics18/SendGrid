package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math/rand"
	"time"
	"bytes"
)

// *** Randomly Generate Data ***

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func MakeRandomEmail() string {
	user := RandStringRunes(10)
	domain := RandStringRunes(5)
	tld := RandStringRunes(3)
	email := fmt.Sprintf("%s@%s.%s", user, domain, tld)
	return email
}

func MakeRandomEmails(n int) []string {
	email_list := make([]string, n)
	for i := range email_list {
		email_list[i] = MakeRandomEmail()
	}
	return email_list
}

func MakeRandomUserID(n int, offset int) int {
	return rand.Intn(n) + offset
}

func MakeRandomUsers(NumOfUsers int, NumOfEmails int) []User {
	ListOfUsers := make([]User, NumOfUsers)
	for i := range ListOfUsers {
		ListOfUsers[i] = User{
			UserID: i,
			Emails: MakeRandomEmails(NumOfEmails),
		}
	}
	return ListOfUsers
}

// **************

// *** Structs *** 

type User struct {
	UserID int 
	Emails []string 
}

type Payload struct {
	UserData User 
}

// **************

func main() {

	numEmails := 10
	userID := 5
	userEmails:= MakeRandomEmails(numEmails)

	d := User{userID, userEmails}
	p := Payload{d}

	userJSON, err := json.MarshalIndent(p, "", "  ")
	req, err := http.NewRequest("GET", "http://localhost:8081/retrieve", bytes.NewBuffer(userJSON))

	req.Header.Set("Content-Type", "application/json")

	// *** Server sends information back to Client *** 

	client := &http.Client{}
	resp, err := client.Do(req)
	
	if err != nil {
		panic(err)
	}

	// Print the response being sent to the client
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response: ", string(body))
	
	resp.Body.Close()
}
