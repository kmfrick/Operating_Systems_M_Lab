package main

import (
	"fmt"
	"math/rand"
	"time"
)


func consume(ch chan int, done chan int, name string) {
	const maxRequests = 5
	const maxConsumeTime = 4
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for cnt := 0; cnt < maxRequests; cnt++ {
		request := <-ch
		fmt.Printf("Consumer %s serving request code %d\n", name, request)
		time.Sleep(time.Second * time.Duration(r.Int()%maxConsumeTime))
		fmt.Printf("Consumer %s served request code %d, num %d\n", name, request, cnt)
	}
	fmt.Printf("Consumer done\n")
	done <- 1
}

const numCodes = 4 // 0 red 1 yellow 2 green 3 white

func produce(ch [numCodes]chan int, name string) {
	const maxProduceTime = 10
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		request := r.Int()%numCodes
		ch[request] <- request
		fmt.Printf("Producer %s sent request code %d\n", name, request)
		time.Sleep(time.Second * time.Duration(r.Int()%maxProduceTime))
	}
}

func serve(triageQueue [numCodes]chan int, triage chan int, name string) {
	for {
		for i := 0; i < numCodes; i++ {
			if len(triageQueue[i]) > 0 {
				request := <-triageQueue[i]
				fmt.Printf("Sending request code %d to consumer %s\n", i, name)
				triage <- request
				break
			}
		}
	}
	panic("Server dying\n")
}


func main() {

	const numAdultConsumers = 10
	const numChildrenConsumers = 10
	const numAdultProducers = 10
	const numChildrenProducers = 10
	const bufSz = 50

	var adultsTriageQueue [numCodes]chan int
	var childrenTriageQueue [numCodes]chan int
	adultsTriage := make(chan int)
	childrenTriage := make(chan int)

	for i := 0; i < numCodes; i++ {
		adultsTriageQueue[i] = make(chan int, bufSz)
		childrenTriageQueue[i] = make(chan int, bufSz)
	}

	var doneAdults [numAdultConsumers]chan int
	var doneChildren [numChildrenConsumers]chan int


	for i := 0; i < numAdultProducers; i++ {
		go produce(adultsTriageQueue, fmt.Sprintf("adults%d", i))
	}

	for i := 0; i < numChildrenProducers; i++ {
		go produce(childrenTriageQueue, fmt.Sprintf("children%d", i))
	}

	go serve(adultsTriageQueue, adultsTriage, "adults")
	go serve(childrenTriageQueue, childrenTriage, "children")

	for i := 0; i < numAdultConsumers; i++ {
		doneAdults[i] = make(chan int)
		go consume(adultsTriage, doneAdults[i], fmt.Sprintf("adults%d", i))
	}
	for i := 0; i < numChildrenConsumers; i++ {
		doneChildren[i] = make(chan int)
		go consume(childrenTriage, doneChildren[i], fmt.Sprintf("children%d", i))
	}

	for i := 0; i < numAdultConsumers; i++ {
		<-doneAdults[i]
	}

	for i := 0; i < numChildrenProducers; i++ {
		<-doneChildren[i]
	}
}



