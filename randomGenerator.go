package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type User struct {
	Userid string   `json: userid`
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
	fmt.Printf("%+v\n", user_list)
	fmt.Println("Done Building")
	userJson, err := json.Marshal(user_list)
	req, err := http.NewRequest("GET", "http://localhost:8081/add", bytes.NewBuffer(userJson))
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
			Userid: strconv.Itoa(i),
			Emails: MakeRandomEmails(NumOfEmails),
		}
	}
	return ListOfUsers
}
