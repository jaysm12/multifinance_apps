package main

import (
	"os"

	"github.com/jaysm12/multifinance-apps/cmd/multifinance-apps/server"
)

func main() {
	os.Exit(server.Run())
}
