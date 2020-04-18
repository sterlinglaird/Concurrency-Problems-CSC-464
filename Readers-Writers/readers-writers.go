package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//type empty {}

type semaphore chan bool

func (s semaphore) wait() {
	s <- true
}

func (s semaphore) notify() {
	<-s
}

var NUM_READERS = 10
var NUM_WRITERS = 10
var NUM_ACTIONS = 100

var wg sync.WaitGroup

var semEmpty = make(semaphore, 1)
var mutex = sync.Mutex{}

var value = 0
var numReaders = 0

func writer() {
	for i := 0; i < NUM_ACTIONS; i++ {
		semEmpty.wait()
		value = rand.Int()
		semEmpty.notify()
	}

	wg.Done()
}

func reader() {
	for i := 0; i < NUM_ACTIONS; i++ {
		mutex.Lock()
		numReaders++
		if numReaders == 1 {
			semEmpty.wait()
		}
		mutex.Unlock()

		//Read value

		mutex.Lock()
		numReaders--
		if numReaders == 0 {
			semEmpty.notify()
		}

		mutex.Unlock()

	}

	wg.Done()
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	NUM_READERS, _ = strconv.Atoi(args[0])
	NUM_WRITERS, _ = strconv.Atoi(args[1])
	NUM_ACTIONS, _ = strconv.Atoi(args[2])

	//Wait for all to finish
	wg.Add(NUM_READERS + NUM_WRITERS)

	for i := 0; i < NUM_READERS; i++ {
		go reader()
	}

	for i := 0; i < NUM_WRITERS; i++ {
		go writer()
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("# readers: %d\n", NUM_READERS)
	fmt.Printf("# writers: %d\n", NUM_WRITERS)
	fmt.Printf("# actions: %d\n", NUM_ACTIONS)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
