package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var NUM_SAVAGES = 10
var MAX_SERVINGS = 5
var NUM_COOK_ITERS = 10

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type Pot struct {
	filledAlert chan bool
	servings    chan int
	waiting     chan bool
	alertNum    int
}

func newPot(numServings int) *Pot {
	pot := new(Pot)

	pot.servings = make(chan int, numServings)
	pot.filledAlert = make(chan bool)
	pot.waiting = make(chan bool)

	return pot
}

func (pot *Pot) getServing() int {
	var serving int

	//If there is no available serving, then wait until the cook fills the pot
	select {
	case serving = <-pot.servings:
		return serving
	default:
		pot.waiting <- true
		<-pot.filledAlert
		serving = pot.getServing()
	}

	return serving
}

func (pot *Pot) fillPot() {
	//Make sure all savages are waiting
	for i := 0; i < NUM_SAVAGES; i++ {
		<-pot.waiting
	}

	//Number all of the servings
	for i := 0; i < MAX_SERVINGS; i++ {
		pot.servings <- i
	}

	//Alert the tribe that the pot is full
	for i := 0; i < NUM_SAVAGES; i++ {
		pot.filledAlert <- true
	}
}

var pot *Pot
var wakeCook = make(chan bool)
var wg sync.WaitGroup

func eat() {
}

func cook() {
	for i := 0; i < NUM_COOK_ITERS; i++ {
		<-wakeCook
		pot.fillPot()
	}

	wg.Done()
}

func savage(id int) {
	for {
		serving := pot.getServing()

		//Alert the cook that the pot is empty if we got the last serving
		if serving == MAX_SERVINGS-1 {
			wakeCook <- true
		}

		eat()
	}
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	NUM_SAVAGES, _ = strconv.Atoi(args[0])
	MAX_SERVINGS, _ = strconv.Atoi(args[1])
	NUM_COOK_ITERS, _ = strconv.Atoi(args[2])

	wg.Add(1) //For the cook

	pot = newPot(MAX_SERVINGS)

	go cook()

	//Get cook to fill pot initially
	wakeCook <- true

	for i := 0; i < NUM_SAVAGES; i++ {
		go savage(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("# savages: %d\n", NUM_SAVAGES)
	fmt.Printf("Max servings: %d\n", MAX_SERVINGS)
	fmt.Printf("# cook iterations: %d\n", NUM_COOK_ITERS)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
