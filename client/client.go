package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	//nolint:gci
	"go-http-server-with-cache/consts"
)

type Client struct {
	id int
	ch chan bool
	wg *sync.WaitGroup
}

func NewClient(id int, wg *sync.WaitGroup) *Client {
	ch := make(chan bool)

	return &Client{id: id, ch: ch, wg: wg}
}

func (c *Client) Run() {
	go func() {
		for {
			select {
			case <-c.ch:
				c.wg.Done()
				log.Println("stopped client ", c.id)

				return
			default:
				//nolint: noctx
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d?clientId=%d", consts.Port, c.id))
				if err != nil {
					log.Printf("Error in client %d, err: %v", c.id, err)
				} else {
					log.Printf("Client %d: %s", c.id, resp.Status)
				}

				resp.Body.Close()

				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			}
		}
	}()
}

func (c *Client) StopClient() {
	c.ch <- true
}
