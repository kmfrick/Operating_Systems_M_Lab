package main

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	name string
	ack  chan int
}

func worker(ch chan req, ack chan int, name string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxSleepTime = 10
	for {
		time.Sleep(time.Duration(r.Int()%maxSleepTime) * time.Second)
		ch <- req{name, ack}
		fmt.Printf("%s sent request, waiting for ack\n", name)
		<-ack
		fmt.Printf("%s received ack\n", name)
	}
}

func when(cond bool, ch chan req) chan req {
	if cond {
		return ch
	}
	return nil
}

func main() {
	const numBuyWorkers = 6
	const numSellWorkers = 6
	const maxFull = 30
	const maxEmpty = 20
	const x = 1
	const y = 1
	const z = 1
	const k = 45

	buyWorker := make(chan req, numBuyWorkers)
	buyWorkerCard := make(chan req, numBuyWorkers)
	sellWorker := make(chan req, numSellWorkers)
	sellWorkerCard := make(chan req, numSellWorkers)
	full := 0
	empty := 0
	cash := 50
	account := 50

	for i := 0; i < numBuyWorkers; i++ {
		buyWorkerAck := make(chan int)
		go worker(buyWorker, buyWorkerAck, "BuyWorker")
	}
	for i := 0; i < numSellWorkers; i++ {
		sellWorkerAck := make(chan int)
		go worker(sellWorker, sellWorkerAck, "SellWorker")
	}
	for i := 0; i < numBuyWorkers; i++ {
		buyWorkerAck := make(chan int)
		go worker(buyWorkerCard, buyWorkerAck, "BuyWorkerCard")
	}
	for i := 0; i < numSellWorkers; i++ {
		sellWorkerAck := make(chan int)
		go worker(sellWorkerCard, sellWorkerAck, "SellWorkerCard")
	}

	done := make(chan int)
	const requestsToServe = 100
	go func() {
		for i := 0; i < requestsToServe; {
			canSell := full >= y && empty+z <= maxEmpty
			canBuy := full+x <= maxFull
			canBuyCash := cash >= z && canBuy
			canBuyCard := account >= z && canBuy
			shouldBuyCash := cash < k || len(buyWorkerCard) == 0 || !canBuyCard
			shouldSellCash := cash >= k || len(sellWorkerCard) == 0
			select {
			case request := <-when(canSell && !shouldSellCash, sellWorkerCard):
				full -= y
				empty += z
				account += z
				request.ack <- 1
				i++
			case request := <-when(canSell && shouldSellCash, sellWorker):
				full -= y
				empty += z
				cash += z
				request.ack <- 1
				i++
			case request := <-when(canBuyCard && !shouldBuyCash, buyWorkerCard):
				empty = 0
				full += x
				account -= x
				request.ack <- 1
				i++
			case request := <-when(canBuyCash && shouldBuyCash, buyWorker):
				empty = 0
				full += x
				cash -= x
				request.ack <- 1
				i++
			}
		}
		done <- 1
	}()

	<-done
}
