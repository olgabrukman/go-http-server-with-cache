package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"go-http-server-with-cache/cache"
	"go-http-server-with-cache/client"
	"go-http-server-with-cache/consts"
	"go-http-server-with-cache/server"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	wg := &sync.WaitGroup{}
	wg.Add(1)

	cache := *cache.NewCache(consts.Max)

	srv := server.NewServer(wg, consts.Port, cache)

	fmt.Println("Enter number of clients")

	var clients int

	_, err := fmt.Scanf("%d", &clients)
	if err != nil {
		log.Panicf("failed to read number of clients, err: %v, quitting", err)
	}

	clientsArr := make([]*client.Client, clients)

	var clientsWg sync.WaitGroup

	for i := 0; i < clients; i++ {
		clientsWg.Add(1)
		clientsArr[i] = client.NewClient(i+1, &clientsWg)
		clientsArr[i].Run()
	}

	fmt.Println("Press the Enter Key to stop server and clients anytime")
	fmt.Scanln()

	log.Println("Stopping clients")

	for i := 0; i < clients; i++ {
		fmt.Println("Sending stop to client ", i+1)
		clientsArr[i].StopClient()
	}

	clientsWg.Wait()
	log.Println("All clients are done")

	srv.StopServer()
}
