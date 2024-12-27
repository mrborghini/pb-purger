package config

import (
	"encoding/json"
	"os"
)

type Collection struct {
	Name                string `json:"name"`
	DeletionTimeSeconds int    `json:"deletionTimeSeconds"`
}

type Entry struct {
	PBUrl               string       `json:"pbUrl"`
	Collections         []Collection `json:"collections"`
	DeletionTimeSeconds int          `json:"deletionTimeSeconds"`
	PBUsername          string       `json:"pbUsername"`
	PBPassword          string       `json:"pbPassword"`
	AccountCollection   string       `json:"accountCollection"`
}

func Read(filePath string) []Entry {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return []Entry{}
	}

	var cfg []Entry
	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return []Entry{}
	}

	return cfg
}

func GetLowestSleepingTime(entries []Entry) int {
	lowestSleepingTime := 0
	for _, entry := range entries {
		for _, collection := range entry.Collections {
			if lowestSleepingTime == 0 || collection.DeletionTimeSeconds < lowestSleepingTime {
				lowestSleepingTime = collection.DeletionTimeSeconds
			}
		}
	}
	return lowestSleepingTime
}
