package main

import (
	"fmt"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type req struct {
	amount int
	ack    chan bool
}


func worker(ch chan req, max int, ack chan bool, name string) {
	const maxSleepTime = 5
	for {
		amount := r.Int()%max + 1
		ch <- req{amount, ack}
		fmt.Printf("%s sent request %d, waiting for ack\n", name, amount)
		<-ack
		fmt.Printf("%s received ack\n", name)
		time.Sleep(time.Duration(r.Int()%maxSleepTime) * time.Second)
	}
}

func serveRequest(ch chan req, subtractFrom *int, addTo *int, max int, mustSubtract bool) bool {
	// Search for prioritary request
	request := req{-1, nil}
	numReq := len(ch)
	for i := 0; i < numReq; i++ {
		request = <-ch
		if request.amount > max/2 || (mustSubtract && request.amount > *subtractFrom) {
//			fmt.Printf("Request %d is not prioritary or not admissible\n", request.amount)
			ch <- request
			request = req{-1, nil}
		} else {
			fmt.Printf("Request %d is prioritary, serving\n", request.amount)
			break // Serve the first prioritary request found
		}
	}
	if request.amount < 0 { // No prioritary request
		numReq := len(ch)
		for i := 0; i < numReq; i++ {
			request = <-ch
			if (*addTo+request.amount > max) || (mustSubtract && request.amount > *subtractFrom) {
//				fmt.Printf("Request %d is not admissible\n", request.amount)
				ch <- request
				request = req{-1, nil}
			} else {
				fmt.Printf("Request %d is admissible, serving\n", request.amount)
				break // Serve the first admissible request found
			}
		}
	}
	if request.amount > 0 { // Serve iff there is an admissible request
		fmt.Printf("Serving request %d\n", request.amount)
		(*subtractFrom) -= request.amount
		(*addTo) += request.amount
		request.ack <- true
		return true
	}
	return false
}

func main() {
	const numHospitalWorkers = 4
	const numLaundryWorkers = 5
	const maxClean = 30
	const maxDirty = 20

	hospitalWorker := make(chan req, numHospitalWorkers)
	laundryWorker := make(chan req, numLaundryWorkers)
	clean := 0
	dirty := 0

	for i := 0; i < numHospitalWorkers; i++ {
		hospitalWorkerAck := make(chan bool)
		go worker(hospitalWorker, maxDirty, hospitalWorkerAck, "HospitalWorker")
	}
	for i := 0; i < numLaundryWorkers; i++ {
		laundryWorkerAck := make(chan bool)
		go worker(laundryWorker, maxClean, laundryWorkerAck, "LaundryWorker")
	}

	for {
		if (len(hospitalWorker)+len(laundryWorker) > 0) {
			hospitalWorkerFirst := clean > dirty
			servedPrioritary := false
			if hospitalWorkerFirst && len(hospitalWorker) > 0 {
				if serveRequest(hospitalWorker, &clean, &dirty, maxDirty, true) {
					fmt.Printf("Clean = %d; Dirty = %d\n", clean, dirty)
					servedPrioritary = true
				} else {
//					fmt.Printf("Could not serve prioritary, falling through\n")
				}
			} 
			if len(laundryWorker) > 0 && !servedPrioritary {
				if serveRequest(laundryWorker, &dirty, &clean, maxClean, false) {
					if dirty < 0 {
						dirty = 0
					}
					fmt.Printf("Clean = %d; Dirty = %d\n", clean, dirty)
				} else {
//					fmt.Printf("Could not serve nonprioritary request\n")
				}
			}
		}
	}
}
