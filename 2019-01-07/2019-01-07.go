package main

import (
	"fmt"
	"math/rand"
	"time"
)

func when(cond bool, ch chan int) chan int {
	if cond {
		return ch
	}
	return nil
}

func main() {
	biscuit := make(chan int)
	biscuitAck := make(chan int)
	iceCream := make(chan int)
	iceCreamAck := make(chan int)
	assembly := make(chan int)
	assemblyAck := make(chan int)
	done := make(chan int)

	go func() { // biscuit production worker
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		const maxSleepTime = 5
		for {
			fmt.Printf("Producing biscuit\n")
			time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
			fmt.Printf("Sending biscuit\n")
			biscuit <- 1
			<-biscuitAck
		}
	}()

	go func() { // ice cream production worker
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		const maxSleepTime = 5
		for {
			<-iceCream
			fmt.Printf("Replenishing ice cream \n")
			time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
			fmt.Printf("Replenished ice cream \n")
			iceCreamAck <- 1
		}
	}()

	go func() { // assembly worker
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		const target = 5
		const maxSleepTime = 5
		for i := 0; i < target; i++ {
			fmt.Printf("Requesting ingredients to assemble ice cream %v\n", i)
			assembly <- 1
			<-assemblyAck
			time.Sleep(time.Second * time.Duration(r.Int()%maxSleepTime))
			fmt.Printf("Ice cream %v assembled\n", i)
		}
		done <- 1
	}()

	go func() { // warehouse
		const max = 5
		numIceCream := 0
		numBiscuit := 0
		for {
			shouldProduceBiscuit := numBiscuit < max/2

			select {
			case <-when(shouldProduceBiscuit, biscuit):
				numBiscuit++
				biscuitAck <- 1
			case <-when(!shouldProduceBiscuit, assembly):
				numBiscuit -= 2
				if (numIceCream == 0) {
					iceCream <- 1
					<-iceCreamAck
					numIceCream = max * 2
				}
				numIceCream--
				assemblyAck <- 1
			}
			fmt.Printf("numIceCream = %v; numBiscuit = %v\n", numIceCream, numBiscuit)
		}
	}()

	<-done
}


