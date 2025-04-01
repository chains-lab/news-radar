package main

import (
	"os"

	"github.com/hs-zavet/news-radar/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
