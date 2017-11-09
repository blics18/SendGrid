package main

import (
	"github.com/sheilatruong96/SendGrid/client"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	//client.PopulateDB()
	 // client.Populate()
	 // client.PopulateDB()

	b := []string{"DxYYwQUtLZ@yahoo.com"}
	client.Check(5, b)


	fmt.Println("Did it work?")
}
