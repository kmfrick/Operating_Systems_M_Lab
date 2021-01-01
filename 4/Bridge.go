package main

import (
	"fmt"
	"math/rand"
	"time"
)

func When(cond bool, ch chan int) chan int {
	if !cond {
		return nil
	}
	return ch
}

func Serve(inN chan int, inS chan int, outN chan int, outS chan int, ack []chan bool) {
	const Capacity = 4
	onBridge := 0
	direction := -1 // 0 N->S 1 S->N
	for {
		if len(outN) > 0 && len(outS) > 0 {
			panic("Road accident!\n")
		}
		canExitS := direction == 0 && onBridge > 0
		canExitN := direction == 1 && onBridge > 0
		canEnterN := onBridge == 0 || (onBridge <= Capacity && direction == 0)
		canEnterS := (onBridge == 0 && len(inN) == 0) || (onBridge <= Capacity && direction == 1)
		select {
		case id := <-When(canExitS, outS):
			fmt.Printf("Vehicle %d exiting from the South\n", id)
			onBridge--
			ack[id] <- true
		case id := <-When(canExitN, outN):
			fmt.Printf("Vehicle %d exiting from the North\n", id)
			onBridge--
			ack[id] <- true
		case id := <-When(canEnterS, inS):
			fmt.Printf("Vehicle %d entering from the South\n", id)
			onBridge++
			direction = 1
			ack[id] <- true
		case id := <-When(canEnterN, inN):
			fmt.Printf("Vehicle %d entering from the North\n", id)
			onBridge++
			direction = 0
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

func Client(id int, chIn chan int, chOut chan int, chAck chan bool, chDone chan bool) {
	const MaxDuration = 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Int()%MaxDuration) * time.Second)
	fmt.Printf("Vehicle %d trying to enter\n", id)
	chIn <- id
	<-chAck
	time.Sleep(time.Duration(r.Int()%MaxDuration) * time.Second)
	fmt.Printf("Vehicle %d trying to exit\n", id)
	chOut <- id
	<-chAck
	chDone <- true
}

func main() {
	const vehiclesNorth = 4
	const vehiclesSouth = 3
	inN := make(chan int, vehiclesNorth)
	inS := make(chan int, vehiclesSouth)
	outN := make(chan int)
	outS := make(chan int)
	var ack = make([]chan bool, vehiclesNorth+vehiclesSouth)
	for i := 0; i < vehiclesNorth+vehiclesSouth; i++ {
		ack[i] = make(chan bool)
	}
	var done [vehiclesNorth + vehiclesSouth]chan bool
	for i := 0; i < vehiclesNorth+vehiclesSouth; i++ {
		done[i] = make(chan bool)
	}
	for i := 0; i < vehiclesNorth; i++ {
		go Client(i, inN, outS, ack[i], done[i])
	}
	for i := vehiclesNorth; i < vehiclesNorth+vehiclesSouth; i++ {
		go Client(i, inS, outN, ack[i], done[i])
	}
	go Serve(inN, inS, outN, outS, ack)
	for i := 0; i < vehiclesNorth+vehiclesSouth; i++ {
		<-done[i]
	}
}
