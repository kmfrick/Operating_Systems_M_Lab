package main

import (
	"fmt"
	"math/rand"
	"time"
)


// Serve takes a message from chIn and moves it one step ahead in the pipeline
// by sending it through chOut
func Serve(stage int, chIn chan int, chOut chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		msg := <-chIn
		fmt.Printf("Stage %d received msg %d.\n", stage, msg)
		t := time.Duration(r.Int() % MaxStageDelay)
		time.Sleep(time.Second * t)
		chOut <- msg
	}
}

// Produce outputs MsgPerProd random numbers on chOut spaced by random time intervals
func Produce(chOut chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const MaxStageDelay = 3
	const MaxMsg = 10
	const MsgPerProd = 3
	for i := 0; i < MsgPerProd; i++ {
		t := time.Duration(r.Int() % MaxStageDelay)
		time.Sleep(time.Second * t)
		msg := r.Int() % MaxMsg

		fmt.Printf("Producing message %d.\n", msg)
		chOut <- (msg)
	}
	fmt.Printf("Producer: done.\n")
	return
}

// Consume infinitely reads from chIn, outputs the message received, and sends a
// "done" signal on chDone
func Consume(chIn chan int, chDone chan bool) {
	for {
		msg := <-chIn
		fmt.Printf("Consuming message %d\n", msg)
		chDone <- true
	}
}

func main() {
	const Stages = 5
	const NumProd = 3
	const NumCons = 4
	producedData := make(chan int)
	consumedData := make(chan int)
	consDone := make(chan bool)
	var bufIn [Stages]chan int
	for i := 0; i < Stages; i++ {
		bufIn[i] = make(chan int)
	}

	// First server reads from producer
	go Serve(0, producedData, bufIn[0])

	for i := 1; i < Stages; i++ {
		go Serve(i, bufIn[i-1], bufIn[i])
	}
	// Last server outputs to consumer
	go Serve(Stages, bufIn[Stages-1], consumedData)

	// Start consumers
	for i := 0; i < NumCons; i++ {
		go Consume(consumedData, consDone)
	}

	// Start producers
	for i := 0; i < NumProd; i++ {
		go Produce(producedData)
	}

	// This could have been done with a global variabe, but is anti-idiomatic
	// and presents data races
	for consumedCnt := 0; consumedCnt < NumProd*MsgPerProd; consumedCnt++ {
		fmt.Printf("Consumed %d messages\n", consumedCnt)
		<-consDone
	}

}
