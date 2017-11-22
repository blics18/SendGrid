package main

import (
	"fmt"

	"github.com/blics18/SendGrid/client"
)

func main() {
	// runClient.go is where we run all of our functions that we write in Client.go.
	// If you'd like to test DropTables and Check, just uncomment the functions.
	// If we decide to make a command-line interface, we'll write it here as well.

	fmt.Println("Starting Client")

	// populate the MySQL Database - numUsers, numEmails, numTables
	db := client.PopulateDB(1, 1, 1)
	defer db.Close()

	// populate the Bloom Filter from values in the MySQL Database
	// client.Populate()

	// to check if values are in the Bloom Filter. Note: Remember to replace the value of b and the userID in Check to what they are in the MySQL DB.
	// b := []string{"NmNTsOQJOl@aol.com"}
	// client.Check(5, b)

	// drop all of the tables in UserStructs schema
	// err := client.DropTables(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//client.HealthCheck()
}
