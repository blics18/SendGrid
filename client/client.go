package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// *** STRUCTS ***
type User struct {
	UserID *int
	Email  []string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Check(userID int, emails []string) error {
	user := User{
		UserID: &userID,
		Email:  emails,
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", "http://localhost:8082/checkBF", bytes.NewBuffer(userJSON))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}

func Clear() error {
	req, err := http.NewRequest("GET", "http://localhost:8082/clearBF", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}

func Populate() error {

	numEmails := 10
	numUsers := 5
	p := MakeRandomUsers(numUsers, numEmails)

	userJSON, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "http://localhost:8082/populateBF", bytes.NewBuffer(userJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	fmt.Println("Response: ", string(body))

	resp.Body.Close()

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeRandomEmail() string {
	user := randStringRunes(10)
	domain := makeEmailDomain()
	email := fmt.Sprintf("%s@%s", user, domain)
	return email
}

func makeRandomEmails(n int) []string {
	email_list := make([]string, n)
	for i := range email_list {
		email_list[i] = makeRandomEmail()
	}
	return email_list
}

func makeEmailDomain() string {
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
		temp := i + 1
		ListOfUsers[i] = User{
			UserID: &(temp), //&i jose
			Email:  makeRandomEmails(NumOfEmails),
		}
	}
	return ListOfUsers
}
