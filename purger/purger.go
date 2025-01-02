package purger

import (
	"fmt"
	"os"
	"pb-purger/config"
	"pb-purger/pb"
	"pb-purger/utils"
	"strings"
	"time"
)

type Purger struct {
	log          *utils.Logger
	sleepingTime int64
}

func Start() {
	purger := newPurger()
	for {
		purger.Run()
	}
}

func newPurger() *Purger {
	return &Purger{
		log:          utils.NewLogger("Purger"),
		sleepingTime: 0,
	}
}

func (p *Purger) checkSleepingTime(timeLeft int64) {
	if p.sleepingTime == 0 || timeLeft < p.sleepingTime {
		p.sleepingTime = timeLeft
	}
}

func (p *Purger) handleUpdatedEntries(searched pb.ListSearch, config config.Collection, pb *pb.PB) {
	for _, entry := range searched.Items {
		normalizedTimestamp := strings.Replace(entry.Updated, " ", "T", 1)
		timestamp, err := time.Parse(time.RFC3339, normalizedTimestamp)
		if err != nil {
			p.log.Warning(fmt.Sprintf("Failed to convert updated time to Unix timestamp: %s", err))
			continue
		}
		unix := timestamp.Unix()
		timePassed := time.Now().Unix() - unix

		if timePassed < 0 {
			timePassed = 0
		}

		if timePassed < int64(config.DeletionTimeSeconds) {
			p.checkSleepingTime(int64(config.DeletionTimeSeconds) - timePassed)
			continue
		}
		success := pb.Delete(entry.Id, config.Name)
		if success {
			p.log.Info(fmt.Sprintf("Deleted entry from %s with ID: %s", config.Name, entry.Id))
			continue
		}
		p.log.Error(fmt.Sprintf("Failed to delete entry from %s with ID: %s", config.Name, entry.Id))
	}
}

// Initialize the purger
func (p *Purger) Run() {
	entries := config.Read("config.json")
	lowestSleepingTime := config.GetLowestSleepingTime(entries)

	if lowestSleepingTime == 0 {
		p.log.Error("deletionTimeSeconds must be set to a value greater than 0")
		os.Exit(1)
	}

	for _, entry := range entries {
		pb := pb.NewPB(entry.PBUrl, entry.PBUsername, entry.PBPassword, entry.AccountCollection)
		if pb.Username != "" || pb.Password != "" {
			err := pb.Login()
			if err != nil {
				p.log.Error(fmt.Sprintf("Failed to login to PB: %s", err))
			}
		}

		for _, collection := range entry.Collections {
			updated := pb.RetrieveLastUpdated(collection.Name)
			p.handleUpdatedEntries(updated, collection, pb)
		}
	}

	if p.sleepingTime == 0 || p.sleepingTime > int64(lowestSleepingTime) {
		p.sleepingTime = int64(lowestSleepingTime)
	}
	p.log.Info(fmt.Sprintf("Sleeping for %d seconds", p.sleepingTime))
	time.Sleep(time.Duration(p.sleepingTime) * time.Second)
	p.sleepingTime = 0
}
