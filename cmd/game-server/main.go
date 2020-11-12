package main

import (
	"fmt"
	"os"

	"github.com/kulti/otus_ol_int_tests/internal/app/gameserver"
)

func main() {
	if err := gameserver.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
