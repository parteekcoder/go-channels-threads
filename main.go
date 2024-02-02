package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tkennon/ticker"
)

const (
	writeDataToRedis string = "writeDataToRedis"
)

type Ticker struct {
	channelG1    chan string
	channelG2    chan string
	countTicking map[string]int
	client       *redis.Client
}

func (t *Ticker) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	t.channelG1 = make(chan string)
	t.channelG2 = make(chan string)
	t.countTicking = make(map[string]int)
	t.client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	go t.G1Ticker(wg)
	go t.G2Ticker(wg)
	go t.Handler(wg)
	wg.Wait()
}

func (t *Ticker) WriteToRedis(key string, val string) {
	t.client.Set(t.client.Context(), key, val, 50*time.Hour)
	fmt.Println(t.client.Get(t.client.Context(), key))
}
func (tick *Ticker) G1Ticker(wg *sync.WaitGroup) {

	t := ticker.NewLinear(time.Second, 1)
	t.Start()
	defer t.Stop()

	for {
		select {
		case <-t.C:
			tick.channelG1 <- "G1"
		case data := <-tick.channelG1:
			if data == writeDataToRedis {
				tick.WriteToRedis("G1", time.Now().UTC().String())
			}
			if data == "Done" {
				wg.Done()
			}
		}
	}
}
func (tick *Ticker) G2Ticker(wg *sync.WaitGroup) {
	t := ticker.NewLinear(time.Second, 1)
	t.Start()
	defer t.Stop()
	for {
		select {
		case <-t.C:
			tick.channelG2 <- "G2"
		case data := <-tick.channelG2:
			if data == "Done" {
				wg.Done()
			}

		}
	}
}

func (t *Ticker) Handler(wg *sync.WaitGroup) {
	for {
		select {
		case <-t.channelG1:
			t.countTicking["G1"]++
			if t.countTicking["G1"] == 3 {
				t.channelG1 <- writeDataToRedis
			}
		case <-t.channelG2:
			t.countTicking["G2"]++
			if t.countTicking["G2"] == 7 {
				t.channelG1 <- writeDataToRedis
				t.channelG2 <- "Done"
				t.channelG1 <- "Done"
				wg.Done()
			}
		}
		fmt.Println(t.countTicking)

	}
}

func main() {
	t := Ticker{}
	t.Start()

}
