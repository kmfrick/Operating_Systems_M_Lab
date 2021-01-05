package main

import (
	"fmt"
	"math/rand"
	"time"
)


// serve takes a message from chIn and moves it one step ahead in the pipeline
// by sending it through chOut
func serve(stage int, chIn chan int, chOut chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		msg := <-chIn
		fmt.Printf("Stage %d received msg %d.\n", stage, msg)
		t := time.Duration(r.Int() % maxStageDelay)
		time.Sleep(time.Second * t)
		chOut <- msg
	}
}

// produce outputs msgPerProd random numbers on chOut spaced by random time intervals
func produce(chOut chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxStageDelay = 3
	const maxmsg = 10
	const msgPerProd = 3
	for i := 0; i < msgPerProd; i++ {
		t := time.Duration(r.Int() % maxStageDelay)
		time.Sleep(time.Second * t)
		msg := r.Int() % maxmsg

		fmt.Printf("Producing message %d.\n", msg)
		chOut <- (msg)
	}
	fmt.Printf("producer: done.\n")
	return
}

// consume infinitely reads from chIn, outputs the message received, and sends a
// "done" signal on chDone
func consume(chIn chan int, chDone chan bool) {
	for {
		msg := <-chIn
		fmt.Printf("Consuming message %d\n", msg)
		chDone <- true
	}
}

func main() {
	const stages = 5
	const numProd = 3
	const numCons = 4
	producedData := make(chan int)
	consumedData := make(chan int)
	consDone := make(chan bool)
	var bufIn [stages]chan int
	for i := 0; i < stages; i++ {
		bufIn[i] = make(chan int)
	}

	// First server reads from producer
	go serve(0, producedData, bufIn[0])

	for i := 1; i < stages; i++ {
		go serve(i, bufIn[i-1], bufIn[i])
	}
	// Last server outputs to consumer
	go serve(stages, bufIn[stages-1], consumedData)

	// Start consumers
	for i := 0; i < numCons; i++ {
		go consume(consumedData, consDone)
	}

	// Start producers
	for i := 0; i < numProd; i++ {
		go produce(producedData)
	}

	// This could have been done with a global variabe, but is anti-idiomatic
	// and presents data races
	for consumedCnt := 0; consumedCnt < numProd*msgPerProd; consumedCnt++ {
		fmt.Printf("consumed %d messages\n", consumedCnt)
		<-consDone
	}

}
