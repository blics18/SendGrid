package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
