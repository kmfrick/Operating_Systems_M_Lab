// resourcepool_prio project main.go

package main

import (
	"fmt"
	"time"
)

const maxProc = 10
const maxRes = 3
const maxBuff = 20

var request = make(chan int, maxBuff)
var release = make(chan int, maxBuff)
var resource [maxProc]chan int
var done = make(chan int)
var terminate = make(chan int)

func client(i int) {

	request <- i
	r := <-resource[i]
	fmt.Printf("\n [Client %d] Using resource %d.\n", i, r)
	time.Sleep(time.Second * 2)
	release <- r
	done <- i
}

func server() {

	var avail int = maxRes
	var res, p, i int
	var free [maxRes]bool
	var susp = 0
	var blocked [maxProc]bool

	// init
	for i := 0; i < maxRes; i++ {
		free[i] = true
	}
	for i := 0; i < maxProc; i++ {
		blocked[i] = false
	}

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("New server cycle")
		select {
		case res = <-release:
			if susp == 0 {
				avail++
				free[res] = true
				fmt.Printf("[server] Releasing resource: %d  \n", res)
			} else {
				for i = maxProc - 1; i >= 0 && !blocked[i]; i-- { // Priority to higher indices
				}
				blocked[i] = false
				susp--
				resource[i] <- res
				fmt.Printf("[server] Woke processo %d - Allocating resource %d  \n", i, res)
			}

		case p = <-request:
			if avail > 0 {
				for i = 0; i < maxRes && !free[i]; i++ {
				}
				free[i] = false
				avail--
				resource[p] <- i
				fmt.Printf("[server]  Allocated resource %d to Client %d \n", i, p)
			} else {
				susp++
				fmt.Printf("[server]  Client %d is waiting... \n", p)
				blocked[p] = true
			}
		case <-terminate:
			fmt.Println("END.")
			done <- 1
			return

		}
	}
}

func main() {
	var cli int
	fmt.Printf("\n  How many clients? (max %d)? ", maxProc)
	fmt.Scanf("%d", &cli)
	fmt.Println("clients:", cli)

	for i := 0; i < cli; i++ {
		resource[i] = make(chan int, maxBuff)
	}

	for i := 0; i < cli; i++ {
		go client(i)
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	terminate <- 1
	<-done

}
