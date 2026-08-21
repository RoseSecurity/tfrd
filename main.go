// Package main acts as the entrypoint for the tfrd command.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/The-Infra-Company/tfrd/cmd"
)

func main() {
	// This sets up a go routine to listen for an interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go listenForInterrupt(stopChan)
	cmd.Execute()
}

func listenForInterrupt(stopScan chan os.Signal) {
	<-stopScan
	log.Fatal("Interrupt signal received. Exiting...")
}
