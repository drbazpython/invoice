package main

import (
	"drbaz.com/invoice/cmd/docprocessing"
	"github.com/charmbracelet/log"
)

func main() {
	err := docprocessing.TestPDFCreation()
	if err != nil {
		log.Fatal("Test failed", "error", err)
	}
	log.Info("Test completed successfully!")
}
