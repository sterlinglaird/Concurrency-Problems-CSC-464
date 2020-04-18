package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var PASSENGERS_PER_CAR = 10
var NUM_PASSENGERS = 1000
var NUM_ROUNDS = 100

var boardQueue chan bool
var unboardQueue chan bool

var wg sync.WaitGroup

func car() {
	for i := 0; i < NUM_ROUNDS; i++ {
		for i := 0; i < PASSENGERS_PER_CAR; i++ {
			<-boardQueue
		}

		for i := 0; i < PASSENGERS_PER_CAR; i++ {
			<-unboardQueue
		}
	}

	wg.Done()
}

func passenger() {
	boardQueue <- true
	unboardQueue <- true
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	PASSENGERS_PER_CAR, _ = strconv.Atoi(args[0])
	NUM_PASSENGERS, _ = strconv.Atoi(args[1])
	NUM_ROUNDS, _ = strconv.Atoi(args[2])

	boardQueue = make(chan bool)
	unboardQueue = make(chan bool)

	//Wait for car to finish
	wg.Add(1)

	go car()

	for i := 0; i < NUM_PASSENGERS; i++ {
		go passenger()
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("passengers per car: %d\n", PASSENGERS_PER_CAR)
	fmt.Printf("# passengers: %d\n", NUM_PASSENGERS)
	fmt.Printf("# rounds: %d\n", NUM_ROUNDS)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
