package main

import (
	"fmt"
	"math/rand"
	"time"
)

const DIM = 5 // Stages
const NUM_PROD = 3
const NUM_CONS = 4
const MAX_STAGE_DELAY = 3
const MSG_PER_PROD = 3
const MAX = 10

// Server take
func Serve(stage int, ch_in chan int, ch_out chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		msg := <-ch_in
		fmt.Printf("Stage %d received msg %d.\n", stage, msg)
		t := time.Duration(r.Int() % MAX_STAGE_DELAY)
		time.Sleep(time.Second * t)
		ch_out <- msg
	}
}

// Produce outputs MSG_PER_PROD random numbers on ch_out spaced by random time intervals
func Produce(ch_out chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < MSG_PER_PROD; i++ {
		t := time.Duration(r.Int() % MAX_STAGE_DELAY)
		time.Sleep(time.Second * t)
		msg := r.Int() % MAX

		fmt.Printf("Producing message %d.\n", msg)
		ch_out <- (msg)
	}
	fmt.Printf("Producer: done.\n")
	return
}

// Consume infinitely reads from ch_in, outputs the message received, and sends a
// "done" signal on ch_done
func Consume(ch_in chan int, ch_done chan bool) {
	for {
		msg := <-ch_in
		fmt.Printf("Consuming message %d\n", msg)
		ch_done <- true
	}
}

func main() {
	produced_data := make(chan int)
	consumed_data := make(chan int)
	cons_done := make(chan bool)
	var buf_in [DIM]chan int
	for i := 0; i < DIM; i++ {
		buf_in[i] = make(chan int)
	}

	// First server reads from producer
	go Serve(0, produced_data, buf_in[0])

	for i := 1; i < DIM; i++ {
		go Serve(i, buf_in[i-1], buf_in[i])
	}
	// Last server outputs to consumer
	go Serve(DIM, buf_in[DIM-1], consumed_data)

	// Start consumers
	for i := 0; i < NUM_CONS; i++ {
		go Consume(consumed_data, cons_done)
	}

	// Start producers
	for i := 0; i < NUM_PROD; i++ {
		go Produce(produced_data)
	}

	// This could have been done with a global variabe, but is anti-idiomatic
	// and presents data races
	for consumed_cnt := 0; consumed_cnt < NUM_PROD*MSG_PER_PROD; consumed_cnt++ {
		fmt.Printf("Consumed %d messages\n", consumed_cnt)
		<-cons_done
	}

}
