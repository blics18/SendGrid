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

	cfg := client.GetEnv()

	client.HealthCheck(cfg)

	db := client.PopulateDB(cfg.NumUsers, cfg.NumEmails, cfg.NumTables)
	defer db.Close()

	client.Populate(cfg)

	fileHandle, err := os.Open("data/data.txt")
	defer fileHandle.Close()
	fileScanner := bufio.NewReader(fileHandle)
	totalMisses := 0
	totalEmails := 0
	totalHits := 0
	userMap := make(map[int][]string)
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
		
		line := strings.Split(s, ":")
		id, email := line[0], line[1]
		userID, _ := strconv.Atoi(id)

		_, exists := userMap[userID]

		if exists {
			userMap[userID] = append(userMap[userID], email)
		} else {
			userMap[userID] = []string{email}
		}
	}

	for key, value := range userMap {
		_, resp := client.Check(cfg, key, value)
		numHits := 0
		numMisses := 0
		numEmails := 0
		totalEmails += len(value)
		numHits += resp.Hits
		numMisses += resp.Miss
		totalMisses += numMisses
		numEmails += resp.Total
		totalHits += numHits
		totalEmails += numEmails
		fmt.Println(fmt.Sprintf("Individual Hit Ratio for User %d: ", key), float64(numHits)/float64(numEmails))
		fmt.Println(fmt.Sprintf("Individual Miss Ratio for User %d: ", key), float64(numMisses)/float64(numEmails))
		fmt.Println()
	}
		
	fmt.Println("Total Hits Ratio: ", float64(totalHits)/float64(totalEmails))
	fmt.Println("Total Miss Ratio: ", float64(totalMisses)/float64(totalEmails))
}

	// drop all of the tables in UserStructs schema
	// err := client.DropTables(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }