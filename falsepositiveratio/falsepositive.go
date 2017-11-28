package main

import (
	"fmt"
	// "strconv"

	"github.com/blics18/SendGrid/client"
)

func testRatio() {
	//CreateBloomFilter(size)
	client.PopulateDB(10, 5000, 5)
	client.Populate()
	users := client.MakeRandomUsers(10, 5000)
	numMisses := 0
	totalEmails := 0
	for i := range users {
		// totalEmails += len(users[i].Email)
		_, resp := client.Check(*users[i].UserID, users[i].Email)
		numMisses += resp.Hits
		totalEmails += resp.Total
		// for x := range resp {
		// 	y, _ := strconv.Atoi(string(resp[x]))
		// 	counter += y
		// }
	}
	fmt.Println(totalEmails)
	fmt.Println(float64(numMisses) / float64(totalEmails))

}

func main() {
	testRatio()
}
