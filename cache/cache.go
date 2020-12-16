package cache

import (
	"log"
	"sync"
	"time"

	"go-http-server-with-cache/consts"
)

//nolint: gochecknoglobals, gomnd
var (
	Interval        = time.Duration(int64(time.Second) * int64(consts.Max))
	CleanupInterval = time.Duration(int64(time.Second) * int64(consts.Max) * 2)
)

type Pair struct {
	timestamp time.Time
	count     int
}

func NewPair() *Pair {
	return &Pair{timestamp: time.Now(), count: 1}
}

type Cache struct {
	m    map[int]*Pair
	max  int
	done chan bool
	lock *sync.RWMutex
}

func NewCache(max int) *Cache {
	m := Cache{
		m:    make(map[int]*Pair),
		max:  max,
		done: make(chan bool),
		lock: &sync.RWMutex{},
	}

	m.cleanUp()

	return &m
}

func (mm *Cache) cleanUp() {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	ticker := time.NewTicker(CleanupInterval)

	go func() {
		for {
			select {
			case <-mm.done:
				ticker.Stop()
				log.Println("stopped background cleaned up process")

				return
			case <-ticker.C:
				if len(mm.m) > 0 {
					for key := range mm.m {
						log.Println("cleaning up: ", key, mm.m[key].count)

						if time.Since(mm.m[key].timestamp).Seconds() > Interval.Seconds() {
							delete(mm.m, key)
						}
					}

					log.Println("cleaned up cache")
				}
			}
		}
	}()
}

func (mm *Cache) Increment(key int) bool {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if mm.m[key] == nil {
		mm.m[key] = NewPair()

		return true
	}

	return mm.m[key].update(mm.max)
}

func (p *Pair) update(max int) bool {
	if time.Since(p.timestamp).Seconds() > Interval.Seconds() {
		p.count = 1
		p.timestamp = time.Now()

		return true
	} else if p.count < max {
		p.count++

		return true
	}

	return false
}

func (mm *Cache) StopCleanUp() {
	mm.done <- true
}
