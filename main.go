package main

import (
	"context"
	"log"
	"os"

	saocli "github.com/rangertaha/sao/internal/cli"
)

func main() {
	app := saocli.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Printf("sao: %v", err)
		os.Exit(1)
	}
}
