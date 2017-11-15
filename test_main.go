package main

import (
	"github.com/blics18/SendGrid/client"
)

func main() {

	//	client.PopulateDB()
	//	client.Populate()

	b := []string{"PgRfZjKhPx@msn.com"}
	client.Check(5, b)

}
