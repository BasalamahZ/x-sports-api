package main

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/x-sports/cmd/xsports-api-http/server"
)

func main() {
	godotenv.Load()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		os.Exit(server.Run())
		defer wg.Done()
	}()
	wg.Wait()
}
