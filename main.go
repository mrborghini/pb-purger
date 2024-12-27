package main

import (
	"fmt"
	"pb-purger/config"
	"pb-purger/pb"
	"strings"
	"time"
)

var sleepingTime int64 = 0

func checkSleepingTime(timeLeft int64) {
	if sleepingTime == 0 || timeLeft < sleepingTime {
		sleepingTime = timeLeft
	}
}

func handleUpdatedEntries(searched pb.ListSearch, config config.Collection) {
	for _, entry := range searched.Items {
		normalizedTimestamp := strings.Replace(entry.Updated, " ", "T", 1)
		timestamp, err := time.Parse(time.RFC3339, normalizedTimestamp)
		if err != nil {
			fmt.Printf("Failed to convert updated time to Unix timestamp: %s\n", err)
			continue
		}
		unix := timestamp.Unix()
		timePassed := time.Now().Unix() - unix

		if timePassed < 0 {
			timePassed = 0
		}

		if timePassed < int64(config.DeletionTimeSeconds) {
			checkSleepingTime(int64(config.DeletionTimeSeconds) - timePassed)
			continue
		}
		fmt.Printf("ID: %s, Updated: %d\n", entry.Id, unix)
	}
}

func run() {
	entries := config.Read("config.json")
	for _, entry := range entries {
		pb := pb.NewPB(entry.PBUrl, entry.PBUsername, entry.PBPassword, entry.AccountCollection)
		if pb.Username != "" || pb.Password != "" {
			err := pb.Login()
			if err != nil {
				fmt.Printf("Failed to login to PB: %s", err)
			}
		}

		for _, collection := range entry.Collections {
			updated := pb.RetrieveLastUpdated(collection.Name)
			handleUpdatedEntries(updated, collection)
		}
	}

	if sleepingTime == 0 {
		sleepingTime = 60
	}
	fmt.Printf("Sleeping for %d seconds\n", sleepingTime)
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	sleepingTime = 0
}

func main() {
	for {
		run()
	}
}
