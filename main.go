package main

import (
	"log"
	"ntc-gwsc/wsc"
	"os"
	"os/signal"
)

var dpwsc *wsc.UWSClient
var cswsc *wsc.UWSClient
var htwsc *wsc.UWSClient
var tkwsc *wsc.UWSClient
var rswsc *wsc.UWSClient

func main() {
	log.Println("=================== Begin Main ===================")

	// DPWSClient
	// dpwsc = wsc.NewDPWSClient()
	// defer dpwsc.Close()
	// go dpwsc.StartDPWSClient()

	// // CSWSClient
	// cswsc = wsc.NewCSWSClient()
	// defer cswsc.Close()
	// go cswsc.StartCSWSClient()

	// // HTWSClient
	// htwsc = wsc.NewHTWSClient()
	// defer htwsc.Close()
	// go htwsc.StartHTWSClient()

	// TKWSClient
	tkwsc = wsc.NewTKWSClient()
	defer tkwsc.Close()
	go tkwsc.StartTKWSClient()

	// // ReloadSymbolWSSClient
	// rswsc = wsc.NewRSWSClient()
	// defer rswsc.Close()
	// go rswsc.StartRSWSClient()

	// Hang thread Main.
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("################# End Main #################")
}
