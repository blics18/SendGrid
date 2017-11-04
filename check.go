package unsub

import (
	"bytes"
	"encoding/json"
	"net/http"
)

/*type User struct {
	UserID int
	Email  []string
}*/

func check(userID int, emails []string) {
	user := User{
		UserID: &userID,
		Email:  emails,
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("GET", "http://localhost:8082/checkBF", bytes.NewBuffer(userJSON))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	resp.Body.Close()
}

/*func main() {
	b := []string{"OYTAEvdVTc@hotmail.com", "hi@gmail.com"}
	check(0, b)
}*/
