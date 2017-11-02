/*
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type User struct {
	Userid int      `json: userid`
	Emails []string `json: emails`
}

func main() {
	nEmails := 10
	//	fmt.Println("Testing MakeRandomEmail")
	//	fmt.Println(MakeRandomEmail())
	//	fmt.Println("Testing MakeRandomEmails")
	MakeRandomEmails(nEmails)
	// for i := 0; i < nEmails; i++ {
	// 	fmt.Println(email_list[i])
	// }
	//	fmt.Println("Testing MakeRandomUsers")
	user_list := MakeRandomUsers(1, 10)
	//fmt.Printf("%+v\n", user_list)
	fmt.Println("Done Building")
	userJson, err := json.Marshal(user_list)
	req, err := http.NewRequest("GET", "http://localhost:8081/add", bytes.NewBuffer(userJson))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response: ", string(body))
	resp.Body.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("%s", string(users))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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

//Do i really need this
func MakeRandomUserID(n int, offset int) int {
	return rand.Intn(n) + offset
}

func MakeRandomUsers(NumOfUsers int, NumOfEmails int) []User {
	ListOfUsers := make([]User, NumOfUsers)
	for i := range ListOfUsers {
		ListOfUsers[i] = User{
			Userid: i,
			Emails: MakeRandomEmails(NumOfEmails),
		}
	}
	return ListOfUsers
} */

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

func MakeRandomUsers(NumOfUsers int, NumOfEmails int) []User {
	ListOfUsers := make([]User, NumOfUsers)
	for i := range ListOfUsers {
		ListOfUsers[i] = User{
			UserID: i,
			Email:  MakeRandomEmails(NumOfEmails),
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

// **************

func main() {

	numEmails := 10
	numUsers := 5
	// MakeRandomEmails(numEmails)
	p := MakeRandomUsers(numUsers, numEmails)

	userJSON, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(err)
	}
	//8082 server
	req, err := http.NewRequest("GET", "http://localhost:8081/retrieve", bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
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
