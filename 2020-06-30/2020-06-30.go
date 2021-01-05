package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Constraints:
// * at most max vehicles on the bridge when in state 0 or 2
// * at most one boat at a time
// * no vehicles or boats in opposite directions
// * can only change state when empty
// * priority (high to low): boats, public transport vehicles, private vehicles
// Bool channels: only true for boats, true for public vehicles, false for private vehicles

type req struct {
	id  int
	ack chan int
}

func when(cond bool, ch chan req) chan req {
	if cond {
		return ch
	}
	return nil
}

func serve(boats, privateVehiclesN, privateVehiclesS, publicVehiclesN, publicVehiclesS, exit chan req) {
	const max = 5
	load := 0
	raised := true
	north := true

	for {
		boatsCanPass := load == 0
		publicVehiclesNCanPass := (!raised && north && load < max) || (load == 0 && len(boats) == 0)
		publicVehiclesSCanPass := (!raised && !north && load < max) || (load == 0 && len(boats) == 0)
		publicVehiclesCnt := len(publicVehiclesN) + len(publicVehiclesS)
		privateVehiclesNCanPass := (!raised && north && load < max) || (load == 0 && len(boats)+publicVehiclesCnt == 0)
		privateVehiclesSCanPass := (!raised && !north && load < max) || (load == 0 && len(boats)+publicVehiclesCnt == 0)

		select {
		case request := <-exit:
			fmt.Printf("%v exiting\n", request.id)
			load--
			request.ack <- 1
		case request := <-when(boatsCanPass, boats):
			raised = true
			load = 1
			request.ack <- 1
		case request := <-when(publicVehiclesNCanPass, publicVehiclesN):
			raised = false
			north = true
			load++
			request.ack <- 1
		case request := <-when(publicVehiclesSCanPass, publicVehiclesS):
			raised = false
			north = false
			load++
			request.ack <- 1
		case request := <-when(privateVehiclesNCanPass, privateVehiclesN):
			raised = false
			north = true
			load++
			request.ack <- 1
		case request := <-when(privateVehiclesSCanPass, privateVehiclesS):
			raised = false
			north = false
			load++
			request.ack <- 1
		}
		fmt.Printf("Bridge raised: %v; Bridge north: %v; Bridge load: %v\n", raised, north, load)
	}
}

func produce(id int, ch chan req, ack chan int, exit chan req, done chan int) {
	time.Sleep(time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const maxSleepTime = 5
	const maxReq = 5
	for i := 0; i < maxReq; i++ {
		fmt.Printf("%v sending request\n", id)
		ch <- req{id, ack}
		<-ack
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
		exit <- req{id, ack}
		<-ack
		fmt.Printf("%v exited\n", id)
		time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
	}
	done <- id
}

func main() {
	const bufSz = 50
	boats := make(chan req, 50)
	privateVehiclesN := make(chan req, 50)
	privateVehiclesS := make(chan req, 50)
	publicVehiclesS := make(chan req, 50)
	publicVehiclesN := make(chan req, 50)
	exit := make(chan req, 50)

	numPriN := 3
	numPriS := 6
	numPubN := 7
	numPubS := 3
	numBoat := 6

	var ack []chan int

	j := 0
	done := make(chan int)
	fmt.Printf("Private north: ")
	for i := 0; i < numPriN; i++ {
		ack = append(ack, make(chan int))
		go produce(j, privateVehiclesN, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nPrivate south: ")
	for i := 0; i < numPriS; i++ {
		ack = append(ack, make(chan int))
		go produce(j, privateVehiclesS, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nPublic north: ")
	for i := 0; i < numPubN; i++ {
		ack = append(ack, make(chan int))
		go produce(j, publicVehiclesN, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nPublic south: ")
	for i := 0; i < numPubS; i++ {
		ack = append(ack, make(chan int))
		go produce(j, publicVehiclesS, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\nBoats: ")
	for i := 0; i < numBoat; i++ {
		ack = append(ack, make(chan int))
		go produce(j, boats, ack[j], exit, done)
		fmt.Printf("%v ", j)
		j++
	}
	fmt.Printf("\n")


	go serve(boats, privateVehiclesN, privateVehiclesS, publicVehiclesN, publicVehiclesS, exit)

	for i := 0; i < j; i++ {
		<-done
	}

}
