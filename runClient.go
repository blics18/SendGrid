package main

import (
	"github.com/blics18/SendGrid/client"

	// "fmt"
)

func main() {
	client.PopulateDB()
	client.Populate()

	// fmt.Println("Start")
	// db := client.PopulateDB()
	// defer db.Close()
	// err := client.DropTables(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// b := []string{"NmNTsOQJOl@aol.com"}
	// client.Check(5, b)
}
