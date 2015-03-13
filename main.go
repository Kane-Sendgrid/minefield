package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Kane-Sendgrid/minefield/api"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal(
			"Usage:\n" +
				"start\n" +
				"mount /mount/point\n" +
				"unmount /mount/point\n" +
				"detonate /mount/point\n" +
				"defuse /mount/point\n")
	}

	switch flag.Arg(0) {
	case "start":
		fmt.Println("starting server...")
		httpServer := api.NewHTTPServer()

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, os.Kill)
		go func() {
			for sig := range signalChannel {
				fmt.Println("Got signal:", sig, ", attempting graceful shutdown...")
				httpServer.Shutdown()
			}
		}()

		httpServer.Serve()
		println("stopped...")
	}
}
