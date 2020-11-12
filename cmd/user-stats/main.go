package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kulti/otus_ol_int_tests/internal/app/userstats"
)

func main() {
	app, err := userstats.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Start()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	<-stopCh
	app.Stop()
}
