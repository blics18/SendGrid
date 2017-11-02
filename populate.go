package unsub

import (
	"fmt"
	"math/rand"
	"time"
)

// *** RANDOMLY GENERATE DATA ***

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

// *** STRUCTS ***

type User struct {
	UserID *int
	Email  []string
}

/*func main() {

	numEmails := 10
	numUsers := 5
	p := MakeRandomUsers(numUsers, numEmails)

	userJSON, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", "http://localhost:8082/populateBF", bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	// *** Server sends information back to Client ***

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response: ", string(body))

	resp.Body.Close()
}*/
