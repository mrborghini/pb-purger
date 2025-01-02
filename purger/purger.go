package purger

import (
	"fmt"
	"strings"
	"time"
	"os"
	"pb-purger/config"
	"pb-purger/pb"
)

var sleepingTime int64 = 0

func checkSleepingTime(timeLeft int64) {
	if sleepingTime == 0 || timeLeft < sleepingTime {
		sleepingTime = timeLeft
	}
}

func handleUpdatedEntries(searched pb.ListSearch, config config.Collection, pb *pb.PB) {
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
		success := pb.Delete(entry.Id, config.Name)
		if success {
			fmt.Printf("Deleted entry from %s with ID: %s\n", config.Name, entry.Id)
			continue
		}
		fmt.Printf("Failed to delete entry from %s with ID: %s\n", config.Name, entry.Id)
	}
}

// Initialize the purger
func Run() {
	entries := config.Read("config.json")
	lowestSleepingTime := config.GetLowestSleepingTime(entries)

	if lowestSleepingTime == 0 {
		fmt.Println("deletionTimeSeconds must be set to a value greater than 0")
		os.Exit(1)
	}

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
			handleUpdatedEntries(updated, collection, pb)
		}
	}

	if sleepingTime == 0 || sleepingTime > int64(lowestSleepingTime) {
		sleepingTime = int64(lowestSleepingTime)
	}
	fmt.Printf("Sleeping for %d seconds\n", sleepingTime)
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	sleepingTime = 0
}