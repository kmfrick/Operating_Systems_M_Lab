package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numWheels = 4
const maxSleepTime = 5

func produce(name string, out chan bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int()%maxSleepTime) * time.Second)
		fmt.Printf("%s producing.\n", name)
		out <- true
	}
}

func assemble(name string, pullC chan bool, pullP chan bool, done chan bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		time.Sleep(time.Duration(r.Int()%maxSleepTime) * time.Second)
		for i := 0; i < numWheels; i++ {
			pullC <- true
			pullP <- true
			fmt.Printf("%s consuming.\n", name)
		}
		done <- true
	}
}

func when(cond bool, ch chan bool) chan bool {
	if cond {
		return ch
	}
	return nil
}

func serve(pushCA, pushCB, pushPA, pushPB, pullCA, pullCB, pullPA, pullPB, doneA, doneB, doneSrv chan bool) {
}

func main() {
	pushCA := make(chan bool)
	pushCB := make(chan bool)
	pushPA := make(chan bool)
	pushPB := make(chan bool)
	pullCA := make(chan bool)
	pullCB := make(chan bool)
	pullPA := make(chan bool)
	pullPB := make(chan bool)

	doneA := make(chan bool)
	doneB := make(chan bool)

	go produce("PA", pushPA)
	go produce("CA", pushCA)
	go produce("PB", pushPB)
	go produce("CB", pushCB)

	go assemble("A", pullCA, pullPA, doneA)
	go assemble("B", pullCB, pullPB, doneB)
	
	const target = 10
	const maxP = 5
	const maxC = 6
	cntCA := 0
	cntCB := 0
	cntPA := 0
	cntPB := 0
	cntDoneA := 0
	cntDoneB := 0
	for ; cntDoneA + cntDoneB < target ; {
		canPushCA := (cntCA+cntCB < maxC) && cntDoneA <= cntDoneB
		canPushCB := (cntCA+cntCB < maxC) && cntDoneB < cntDoneA
		canPushPA := (cntPA+cntPB < maxP) && cntDoneA <= cntDoneB
		canPushPB := (cntPA+cntPB < maxP) && cntDoneB < cntDoneA

		canPullCA := cntCA > 0 && (cntDoneA <= cntDoneB || cntCB == 0)
		canPullCB := cntCB > 0 && (cntDoneB < cntDoneA || cntCA == 0)
		canPullPA := cntPA > 0 && (cntDoneA <= cntDoneB || cntPB == 0)
		canPullPB := cntPB > 0 && (cntDoneB < cntDoneA || cntPA == 0)

		select {
		case <-when(canPushCA, pushCA):
			cntCA++
		case <-when(canPushPA, pushPA):
			cntPA++

		case <-when(canPullCA, pullCA):
			cntCA--
		case <-when(canPullPA, pullPA):
			cntPA--

		case <-when(canPullCB, pullCB):
			cntCB--
		case <-when(canPullPB, pullPB):
			cntPB--

		case <-when(canPushCB, pushCB):
			cntCB++
		case <-when(canPushPB, pushPB):
			cntPB++

		case <-doneA:
			cntDoneA++
		case <-doneB:
			cntDoneB++
		}
		fmt.Printf("CA = %d; CB = %d; PA = %d; PB = %d\n", cntCA, cntCB, cntPA, cntPB)
		fmt.Printf("produced %d A cars and %d B cars\n", cntDoneA, cntDoneB)
	}

}
