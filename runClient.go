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
	fmt.Println("Starting Client")
	
	client.HealthCheck()

	for {
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
	}
}

	// drop all of the tables in UserStructs schema
	// err := client.DropTables(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }