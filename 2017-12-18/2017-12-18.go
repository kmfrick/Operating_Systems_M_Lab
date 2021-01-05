package main

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	id  int
	msc bool
	ack chan int
}

func when(cond bool, ch chan req) chan req {
	if cond {
		return ch
	}
	return nil
}

func serve(mscStudGrad, bscStud, mscStud, bscStudGrad, exit chan req) {
	const max = 5
	numBscStud := 0
	numMscStud := 0

	for {
		if numBscStud + numMscStud > max {
			panic("COVID-19 prevention rules are being violated!\n")
		}
		libraryIsFull := numBscStud + numMscStud == max
		mscStudGradHasPriority := !libraryIsFull && numMscStud <= numBscStud
		mscStudHasPriority := !libraryIsFull && mscStudGradHasPriority && len(mscStudGrad) == 0
		bscStudGradHasPriority := !libraryIsFull && !mscStudGradHasPriority
		bscStudHasPriority := !libraryIsFull && bscStudGradHasPriority && len(bscStudGrad) == 0

		select {
		case request := <-exit:
			fmt.Printf("%v exiting\n", request.id)
			if request.msc {
				numMscStud--
			} else {
				numBscStud--
			}
			request.ack <- 1
		case request := <-when(mscStudGradHasPriority, mscStudGrad):
			numMscStud++
			request.ack <- 1
		case request := <-when(mscStudHasPriority, mscStud):
			numMscStud++
			request.ack <- 1
		case request := <-when(bscStudGradHasPriority, bscStudGrad):
			numBscStud++
			request.ack <- 1
		case request := <-when(bscStudHasPriority, bscStud):
			numBscStud++
			request.ack <- 1
		}
		fmt.Printf("BSc: %v; MSc: %v\n", numBscStud, numMscStud)
	}
}

func produce(id int, msc bool, ch chan req, ack chan int, exit chan req, done chan int) {
	time.Sleep(time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxSleepTime = 5
	const maxReq = 5
	for i := 0; i < maxReq; i++ {
		fmt.Printf("%v sending request\n", id)
		ch <- req{id, msc, ack}
		<-ack
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
		exit <- req{id, msc, ack}
		<-ack
		fmt.Printf("%v exited\n", id)
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
	}
	done <- id
}

func main() {
	const bufSz = 50
	mscStudGrad := make(chan req, 50)
	bscStud := make(chan req, 50)
	bscStudGrad := make(chan req, 50)
	mscStud := make(chan req, 50)
	exit := make(chan req, 50)

	numBscStud := 3
	numMscStud := 7
	numBscStudGrad := 3
	numMscStudGrad := 6

	var ack []chan int

	j := 0
	done := make(chan int)
	fmt.Printf("BSc non-grad: ")
	for i := 0; i < numBscStud; i++ {
		ack = append(ack, make(chan int))
		go produce(j, false, bscStud, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nMSc non-grad: ")
	for i := 0; i < numMscStud; i++ {
		ack = append(ack, make(chan int))
		go produce(j, true, mscStud, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nBSc grad: ")
	for i := 0; i < numBscStudGrad; i++ {
		ack = append(ack, make(chan int))
		go produce(j, false, bscStudGrad, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nMSc grad: ")
	for i := 0; i < numMscStudGrad; i++ {
		ack = append(ack, make(chan int))
		go produce(j, true, mscStudGrad, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\n")


	go serve(mscStudGrad, bscStud, mscStud, bscStudGrad, exit)

	for i := 0; i < j; i++ {
		<-done
	}

}
