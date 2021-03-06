package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/defsky/dmon/app"
)

func main() {
	go waitSignal()

	app.Init()
	app.Start(30)
}

func waitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)

	for {
		s := <-c
		switch s {
		case os.Interrupt:
			log.Println("User Interrupt")
			os.Exit(0)
		}
	}
}
