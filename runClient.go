package main

import (
	"fmt"

	"github.com/blics18/SendGrid/client"
)

func main() {
	fmt.Println("Starting Client...")

	cfg := client.GetEnv()

	client.HealthCheck(cfg)

	db, err := client.PopulateDB(cfg.NumUsers, cfg.NumEmails, cfg.NumTables)

	if err != nil {
		fmt.Println("Unable to populate database")
		return
	}

	defer db.Close()

	client.Populate(cfg)

	userMap := client.ParseFile()

	totalMisses := 0
	totalEmails := 0
	totalHits := 0

	for userID, userEmails := range userMap {
		resp, _ := client.Check(cfg, userID, userEmails)

		totalMisses += resp.Miss
		totalHits += resp.Hits
		totalEmails += resp.NumEmails

		fmt.Println(fmt.Sprintf("Individual Hit Ratio for User %d: ", userID), float64(resp.Hits)/float64(resp.NumEmails))
		fmt.Println(fmt.Sprintf("Individual Miss Ratio for User %d: ", userID), float64(resp.Miss)/float64(resp.NumEmails))
		fmt.Println()
	}

	fmt.Println("Total Hits Ratio: ", float64(totalHits)/float64(totalEmails))
	fmt.Println("Total Miss Ratio: ", float64(totalMisses)/float64(totalEmails))
}