package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type semaphore chan bool

func (s semaphore) wait() {
	s <- true
}

func (s semaphore) notify() {
	<-s
}

var NUM_ROUNDS = 10000
var NUM_STUDENTS = 10

var okToLeave = make(semaphore, 0)
var mutex = sync.Mutex{}

var numEating = 0
var numReadyToLeave = 0

var wg sync.WaitGroup

func dine() {
	mutex.Lock()

	numEating++

	if numEating == 2 && numReadyToLeave == 1 {
		okToLeave.notify()
	}

	//Eat!!

	numEating--

	mutex.Unlock()
}

func leave() {
	mutex.Lock()

	numReadyToLeave++

	if numEating == 1 && numReadyToLeave == 1 {
		mutex.Unlock()
		okToLeave.wait()
		mutex.Lock()
		numReadyToLeave--
	} else if numEating == 0 && numReadyToLeave == 2 {
		okToLeave.notify()
		numReadyToLeave--
	} else {
		numReadyToLeave--
	}

	mutex.Unlock()
}

func student() {
	for i := 0; i < NUM_ROUNDS; i++ {
		dine()
		leave()
	}

	wg.Done()
}

func main() {
	start := time.Now()

	args := os.Args[1:]

	NUM_STUDENTS, _ = strconv.Atoi(args[0])
	NUM_ROUNDS, _ = strconv.Atoi(args[1])

	//Wait for car to finish
	wg.Add(NUM_STUDENTS)

	for i := 0; i < NUM_STUDENTS; i++ {
		go student()
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("# students: %d\n", NUM_STUDENTS)
	fmt.Printf("# rounds: %d\n", NUM_ROUNDS)
	fmt.Printf("Time taken: %dms\n", elapsed.Nanoseconds()/1000000)
}
