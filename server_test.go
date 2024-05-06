package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"redos/client"
)

func TestServerWithManyClient(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.start())
	}()
	time.Sleep(time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			defer wg.Done()

			c, err := client.NewClient("localhost:8888")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", it)
			value := fmt.Sprintf("client_bar_%d", it)
			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}

			val, err := c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("client got this value :", val)
		}(i)
	}
	wg.Wait()

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		log.Fatalf("%d peers in peer pool\n", len(server.peers))
	}

}
