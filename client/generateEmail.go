package client

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
	email_list := make([]string, rand.Intn(n)+1)
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

func MakeRandomUsers(NumOfUsers int, MaxNumOfEmails int) []User {
	ListOfUsers := make([]User, NumOfUsers)
	for i := range ListOfUsers {
		temp := i + 1
		ListOfUsers[i] = User{
			UserID: &(temp),
			Email:  makeRandomEmails(MaxNumOfEmails),
		}
	}
	return ListOfUsers
}
