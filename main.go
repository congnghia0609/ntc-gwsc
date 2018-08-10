package main

import (
	"log"
	"ntc-gwsc/wsc"
	"os"
	"os/signal"
)

var tkwsc *wsc.UWSClient

func main() {
	log.Println("=================== Begin Main ===================")

	// TKWSClient
	tkwsc = wsc.NewTKWSClient()
	defer tkwsc.Close()
	go tkwsc.StartTKWSClient()

	// Hang thread Main.
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("################# End Main #################")
}
