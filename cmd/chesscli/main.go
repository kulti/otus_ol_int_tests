package main

import (
	"fmt"
	"os"

	"github.com/kulti/otus_ol_int_tests/internal/app/chesscli"
)

func main() {
	if err := chesscli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
