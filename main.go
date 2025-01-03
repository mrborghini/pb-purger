package main

import (
	"pb-purger/purger"
	"pb-purger/utils"
)

func main() {
	log := utils.NewLogger("Main")
	log.Info("Starting purger")
	// Start the purger
	purger.Start()
}
