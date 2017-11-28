package main

import (
	"github.com/blics18/SendGrid/client"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// runClient.go is where we run all of our functions that we write in Client.go.
	// If you'd like to test DropTables and Check, just uncomment the functions.
	// If we decide to make a command-line interface, we'll write it here as well.

	fmt.Println("Starting Client")

	// populate the MySQL Database - numUsers, numEmails, numTables
	db := client.PopulateDB(5, 1000, 5)
	defer db.Close()

	// populate the Bloom Filter from values in the MySQL Database
	client.Populate()

	fileHandle, err := os.Open("data/data.txt")
	defer fileHandle.Close()
	fileScanner := bufio.NewReader(fileHandle)
	totalMisses := 0
	totalEmails := 0
	for {
		var buffer bytes.Buffer
		var l []byte
		var isPrefix bool
		for {
			l, isPrefix, err = fileScanner.ReadLine()
			buffer.Write(l)
			if !isPrefix {
				break
			}
			if err != nil {
				break
			}
		}
		if err == io.EOF {
			break
		}
		s := buffer.String()
		numMisses := 0
		numEmails := 0
		line := strings.Split(s, ":")
		id, emails := line[0], line[1]
		userID, _ := strconv.Atoi(id)
		userEmails := strings.Split(emails, " ")
		_, resp := client.Check(userID, userEmails)
		totalEmails += len(userEmails)
		numMisses += resp.Hits
		totalMisses += numMisses
		numEmails += resp.Total
		totalEmails += numEmails
		fmt.Println("Individual Ratio: ", float64(numMisses)/float64(numEmails))
	}
	fmt.Println("Total Ratio: ", float64(totalMisses)/float64(totalEmails))
	/*for fileScanner.Scan() {
		s := strings.Split(fileScanner.Text(), ":")
		id, emails := s[0], s[1]
		userID, _ := strconv.Atoi(id)
		userEmails := strings.Split(emails, " ")
		_, resp := client.Check(userID, userEmails)
		numMisses += resp.Hits
		totalEmails += resp.Total
		fmt.Println("Ratio: ", float64(resp.Hits) / float64(resp.Total))
	}*/
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
