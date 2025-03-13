package news_radar

import (
	"os"

	"github.com/recovery-flow/news-radar/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
