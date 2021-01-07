package main

import (
	"fmt"
	"math/rand"
	"time"
)

func produce(id int, exchangeType int, ch chan int, ack chan int, chPrio chan int, ackPrio chan int) {
	const maxSleepTime = 8
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		ch <- id
		fmt.Printf("Producer %v produced request for exchange type %v.\n", id, exchangeType)
		idCons := <-ack
		fmt.Printf("Request %v for exchange type %v was received by consumer %v.\n", id, exchangeType, idCons)
		if ch != chPrio {
			fmt.Printf("Producer %v becoming prioritary\n", id)
			ch = chPrio
			ack = ackPrio
		}
		exchangeType = (exchangeType + 2 + 2*(exchangeType%2) + r.Int()%2) % 6
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
	}
}

func consume(id int, exchangeType int, ch chan int, ack chan int, chPrio chan int, ackPrio chan int) {
	const maxSleepTime = 8
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		ch <- id
		fmt.Printf("Consumer %v waiting to consume exchange type %v.\n", id, exchangeType)
		idProd := <-ack
		fmt.Printf("Consumer %v consumed production %v for exchange type %v.\n", id, idProd, exchangeType)
		if ch != chPrio {
			fmt.Printf("Consumer %v becoming prioritary\n", id)
			ch = chPrio
			ack = ackPrio
		}
		exchangeType = (exchangeType+3 + r.Int()%2) % 6
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
	}
}

func when(cond bool, ch chan int) chan int {
	if cond {
		return ch
	}
	return nil
}

func serve(prodCh chan int, prodAck []chan int, prodChPrio chan int, prodAckPrio []chan int, consCh chan int, consAck []chan int, exchangeType int, done chan int) {
	const requestsToServe = 8
	for i := 0; i < requestsToServe; {
		select {
		case idProd := <-when(len(consCh) > 0, prodChPrio):
			idCons := <-consCh
			consAck[idCons] <- idProd
			prodAckPrio[idProd] <- idCons
			i++
			fmt.Printf("Consumed %v requests for exchange type %v\n", i, exchangeType)
		case idProd := <-when(len(consCh) > 0 && len(prodChPrio) == 0, prodCh):
			idCons := <-consCh
			consAck[idCons] <- idProd
			prodAck[idProd] <- idCons
			i++
			fmt.Printf("Consumed %v requests for exchange type %v\n", i, exchangeType)
		default:
			// noop
		}
	}
	done <- 1
}

func main() {
	const bufSz = 50
	const numProducers = 6
	const numConsumers = 6
	const exchangeTypes = 3 * 2 * 1 // 3 types of bike can be exchanged with one another for 3! types of exchange
	// Man for woman
	// Man for child
	// Woman for child
	// Woman for man
	// Child for man
	// Child for woman
	// Producer's semantics is "giving for accepting", consumer's semantics is "accepting for giving"
	// Next exchange type for a producer is given by: ((type%2 == 0 ? type+2 : type+4) + r.Int()%2) % 6 == (type + 2 + 2*(type%2) + r.Int()%2) % 6
	// Next exchange type for a consumer is given by: (type+3 + r.Int()%2) % 6

	done := make([]chan int, exchangeTypes)

	for k := 0; k < exchangeTypes; k++ {
		consCh := make(chan int, bufSz)
		consChPrio := make(chan int, bufSz)
		consAck := make([]chan int, numConsumers)
		consAckPrio := make([]chan int, numConsumers)

		for i := 0; i < numConsumers; i++ {
			consAck[i] = make(chan int)
			go consume(i, k, consCh, consAck[i], consChPrio, consAckPrio[i])
		}

		prodCh := make(chan int, bufSz)
		prodChPrio := make(chan int, bufSz)
		prodAck := make([]chan int, numProducers)
		prodAckPrio := make([]chan int, numProducers)

		for i := 0; i < numProducers; i++ {
			prodAck[i] = make(chan int)
			prodAckPrio[i] = make(chan int)
			go produce(i, k, prodCh, prodAck[i], prodChPrio, prodAckPrio[i])
		}

		done[k] = make(chan int)
		go serve(prodCh, prodAck, prodChPrio, prodAckPrio, consCh, consAck, k, done[k])
	}
	for k := 0; k < exchangeTypes; k++ {
		<-done[k]
	}

}
