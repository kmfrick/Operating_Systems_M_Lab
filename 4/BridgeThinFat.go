package main

import (
	"fmt"
	"math/rand"
	"time"
)

func when(cond bool, ch chan int) chan int {
	if !cond {
		return nil
	}
	return ch
}

func serve(inFatN chan int, inFatS chan int, inN chan int, inS chan int, outN chan int, outS chan int, ack []chan bool) {
	const capacity = 2
	fatOnBridge := false
	onBridge := 0
	direction := -1 // 0 N->S 1 S->N
	for {
		if len(outN) > 0 && len(outS) > 0 {
			panic("Road accident!\n")
		}
		if (onBridge > capacity || (onBridge > 0 && fatOnBridge)) {
			fmt.Printf("onBridge = %d; fatOnBridge = %d\n", onBridge, fatOnBridge)
			panic("Bridge collapsing!\n")
		}
		canExitS := direction == 0 && (onBridge > 0 || fatOnBridge)
		canExitN := direction == 1 && (onBridge > 0 || fatOnBridge)
		canEnterN := !fatOnBridge && (onBridge == 0 || (onBridge < capacity && direction == 0))
		canEnterS := !fatOnBridge && ((onBridge == 0 && len(inN) == 0) || (onBridge < capacity && direction == 1))
		canEnterFatN := !fatOnBridge && onBridge == 0 && len(inN) == 0
		canEnterFatS := canEnterFatN && len(inS) == 0 && len(inFatN) == 0
		select {
		case id := <-when(canExitS, outS):
			fmt.Printf("Vehicle %d exiting from the South\n", id)
			if (fatOnBridge) {
				fatOnBridge = false
			} else {
				onBridge--
			}
			ack[id] <- true
		case id := <-when(canExitN, outN):
			fmt.Printf("Vehicle %d exiting from the North\n", id)
			if (fatOnBridge) {
				fatOnBridge = false
			} else {
				onBridge--
			}
			ack[id] <- true
		case id := <-when(canEnterS, inS):
			fmt.Printf("Vehicle %d entering from the South\n", id)
			onBridge++
			direction = 1
			ack[id] <- true
		case id := <-when(canEnterN, inN):
			fmt.Printf("Vehicle %d entering from the North\n", id)
			onBridge++
			direction = 0
			ack[id] <- true
		case id := <-when(canEnterFatN, inFatN):
			fmt.Printf("Vehicle %d entering from the North\n", id)
			fmt.Printf("Vehicle %d is fat\n", id)
			fatOnBridge = true
			direction = 0
			ack[id] <- true
		case id := <-when(canEnterFatS, inFatS) :
			fmt.Printf("Vehicle %d entering from the South\n", id)
			fmt.Printf("Vehicle %d is fat\n", id)
			fatOnBridge = true
			direction = 1
			ack[id] <- true
		}
		if direction == 0 {
			fmt.Printf("Direction = ->S\n")
		} else {
			fmt.Printf("Direction = ->N\n")
		}
		fmt.Printf("%d vehicles on the bridge.\n", onBridge)
	}
}

func produce(id int, chIn chan int, chOut chan int, chAck chan bool, chDone chan bool) {
	const maxDuration = 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Int()%maxDuration) * time.Second)
	fmt.Printf("Vehicle %d trying to enter\n", id)
	chIn <- id
	<-chAck
	time.Sleep(time.Duration(r.Int()%maxDuration) * time.Second)
	fmt.Printf("Vehicle %d trying to exit\n", id)
	chOut <- id
	<-chAck
	chDone <- true
}

func main() {
	const vehiclesFatNorth = 2
	const vehiclesFatSouth = 1
	const vehiclesNorth = 4
	const vehiclesSouth = 1
	vehicles := vehiclesNorth+vehiclesSouth+vehiclesFatNorth+vehiclesFatSouth
	inFatN := make(chan int)
	inFatS := make(chan int)
	inN := make(chan int, vehiclesNorth)
	inS := make(chan int, vehiclesSouth)
	outN := make(chan int)
	outS := make(chan int)
	ack := make([]chan bool, vehicles)
	for i := 0; i < vehicles ; i++ {
		ack[i] = make(chan bool)
	}
	done := make([]chan bool, vehicles)
	for i := 0; i < vehicles; i++ {
		done[i] = make(chan bool)
	}
	for i := 0; i < vehiclesNorth; i++ {
		go produce(i, inN, outS, ack[i], done[i])
	}
	for i := vehiclesNorth; i < vehiclesNorth+vehiclesSouth; i++ {
		go produce(i, inS, outN, ack[i], done[i])
	}
	for i := vehiclesNorth+vehiclesSouth; i < vehiclesNorth+vehiclesSouth+vehiclesFatNorth; i++ {
		go produce(i, inFatN, outS, ack[i], done[i])
	}
	for i := vehiclesNorth+vehiclesSouth+vehiclesFatNorth; i < vehicles; i++ {
		go produce(i, inFatS, outN, ack[i], done[i])
	}
	go serve(inFatN, inFatS, inN, inS, outN, outS, ack)
	for i := 0; i < vehicles; i++ {
		<-done[i]
	}
}
