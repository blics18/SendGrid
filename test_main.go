package main

import (
	"SendGrid/client"
	"fmt"
)

func main() {
	// client.Populate()

	fmt.Println("Start")
	db := client.PopulateDB()
	defer db.Close()
	err := client.DropTables(db)
	if err != nil {
		fmt.Println(err)
	}
	//	b := []string{"DxYYwQUtLZ@yahoo.com"}
	//	client.Check(5, b)

}
