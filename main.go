package main

import (
	"sand_clock/sand_clock"
	"sync"
)

func main() {
	stop := make(chan struct{})
	done := sync.WaitGroup{}
	done.Add(1)
	go sand_clock.Drop(stop, &done)
	done.Wait()
}
