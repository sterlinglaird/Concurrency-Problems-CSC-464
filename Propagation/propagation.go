package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var NUM_COPIES = 100
var MAX_DELAY = 50 //ms, use 0 for benchmarking the syncronization
var NUM_ROUNDS = 10

var wg sync.WaitGroup

var ch []chan Modification

type Origin int

const (
	mutation Origin = iota + 1
	propagation
)

type Modification struct {
	newValue int
	origin   Origin
}

func mutator(id int) {

	for i := 0; i < NUM_ROUNDS; i++ {
		if MAX_DELAY != 0 {
			var msToSleep = time.Duration(rand.Int() % MAX_DELAY)
			time.Sleep(msToSleep * time.Millisecond)
		}

		var newValue = rand.Int() % 10000 //Makes number more manageable to look at

		ch[id] <- Modification{newValue: newValue, origin: mutation}
	}

	wg.Done()
}

func propagator(id int) {

	//Only one version mutates
	if id == 0 {
		go mutator(id)
	}

	var value = 0

	for {
		mod := <-ch[id]

		value = mod.newValue

		//Only propagate if it is a mutation
		if mod.origin == mutation {
			for propagateId := 0; propagateId < NUM_COPIES; propagateId++ {
				if propagateId == id {
					continue
				}

				ch[propagateId] <- Modification{newValue: value, origin: propagation}
			}
		}
	}
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	NUM_COPIES, _ = strconv.Atoi(args[0])
	MAX_DELAY, _ = strconv.Atoi(args[1])
	NUM_ROUNDS, _ = strconv.Atoi(args[2])

	//Make all of the resource channels
	ch = make([]chan Modification, NUM_COPIES)

	//Wait for all to finish
	wg.Add(1)

	//open channels, need to do this before we start the propagators
	//or else they might access channels that havent been initialized yet
	for i := 0; i < NUM_COPIES; i++ {
		ch[i] = make(chan Modification)
	}

	for i := 0; i < NUM_COPIES; i++ {
		go propagator(i)
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("# copies: %d\n", NUM_COPIES)
	fmt.Printf("max delay: %dms\n", MAX_DELAY)
	fmt.Printf("# rounds: %d\n", NUM_ROUNDS)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
