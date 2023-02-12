package main

import (
	"server/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.StartServer()
	}()

	wg.Wait()
}
