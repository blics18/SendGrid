package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
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
	domain := MakeEmailDomain()
	email := fmt.Sprintf("%s@%s", user, domain)
	return email
}

func MakeRandomEmails(n int) []string {
	email_list := make([]string, n)
	for i := range email_list {
		email_list[i] = MakeRandomEmail()
	}
	return email_list
}
func MakeEmailDomain() string {
	suffix := make([]string, 0)
	suffix = append(suffix,
		"gmail.com",
		"hotmail.com",
		"yahoo.com",
		"msn.com",
		"aol.com")
	net := suffix[rand.Intn(len(suffix))]
	return net
}

func MakeRandomUsers(NumOfUsers int) []User {
	ListOfUsers := make([]User, NumOfUsers)
	for i := range ListOfUsers {
		ListOfUsers[i] = User{
			UserID: i,
			Email:  MakeRandomEmails(NumOfUsers),
		}
	}
	return ListOfUsers
}

// **************

// *** Structs ***

type User struct {
	UserID int
	Email  []string
}

type Payload struct {
	UserData User
}

// **************

func main() {

	numEmails := 10
	// MakeRandomEmails(numEmails)
	p := MakeRandomUsers(numEmails)

	userJSON, err := json.MarshalIndent(p, "", "  ")
	//8082 server
	req, err := http.NewRequest("GET", "http://localhost:8082/retrieve", bytes.NewBuffer(userJSON))
	fmt.Println(req, err)
	req.Header.Set("Content-Type", "application/json")

	// *** Server sends information back to Client ***

	client := &http.Client{}
	resp, err := client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		panic(err)
	}

	// Print the response being sent to the client
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response: ", string(body))
	// fmt.Println(body)

	resp.Body.Close()
}
