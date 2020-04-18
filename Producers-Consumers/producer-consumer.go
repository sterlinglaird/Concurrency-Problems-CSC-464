package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var BUFFER_SIZE = 100
var NUM_PRODUCERS = 2
var NUM_CONSUMERS = 10
var NUM_PRODUCED = 10000 //each

var buf chan int
var wg sync.WaitGroup

func waitForEvent() int {
	var randomNum = rand.Int()
	return randomNum % 10000
}

func consumeEvent(event int) {

}

func producer(id int) {
	for i := 0; i < NUM_PRODUCED; i++ {
		var event = waitForEvent()
		buf <- event
	}

	wg.Done()
}

func consumer(id int) {
	for {
		event := <-buf
		consumeEvent(event)
	}
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	BUFFER_SIZE, _ = strconv.Atoi(args[0])
	NUM_PRODUCERS, _ = strconv.Atoi(args[1])
	NUM_CONSUMERS, _ = strconv.Atoi(args[2])
	NUM_PRODUCED, _ = strconv.Atoi(args[3])

	buf = make(chan int, BUFFER_SIZE)

	wg.Add(NUM_PRODUCERS)

	for i := 0; i < NUM_PRODUCERS; i++ {
		go producer(i)
	}

	for i := 0; i < NUM_CONSUMERS; i++ {
		go consumer(i)
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("buffer size: %d\n", BUFFER_SIZE)
	fmt.Printf("# producers: %d\n", NUM_PRODUCERS)
	fmt.Printf("# consumers: %d\n", NUM_CONSUMERS)
	fmt.Printf("# produced per producer: %d\n", NUM_PRODUCED)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
