package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("localhost:8888")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	time.Sleep(time.Second)

	for i := 0; i < 5; i++ {
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("data is :", val)
	}

}

func TestClientIntData(t *testing.T) {
	c, err := NewClient("localhost:8888")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	time.Sleep(time.Second)

	if err := c.Set(context.TODO(), "foo", "1"); err != nil {
		log.Fatal(err)
	}

	val, err := c.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)

}
