// resourcepool_prio project main.go

package main

import (
	"fmt"
	"time"
)

const MaxProc = 10
const MaxRes = 3
const MaxBuff = 20
const UsageTime = 2

var request = make(chan int, MaxBuff)
var release = make(chan int, MaxBuff)
var resource [MaxProc]chan int
var done = make(chan int)
var terminate = make(chan int)

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func client(i int) {

	request <- i
	r := <-resource[i]
	fmt.Printf("\n [Client %d] Using resource %d.\n", i, r)
	time.Sleep(time.Second * UsageTime)
	release <- r
	done <- i
}

func server() {

	var avail int = MaxRes
	var res, p, i int
	var free [MaxRes]bool
	var susp = 0
	var blocked [MaxProc]bool

	// init
	for i := 0; i < MaxRes; i++ {
		free[i] = true
	}
	for i := 0; i < MaxProc; i++ {
		blocked[i] = false
	}

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("New server cycle starting...")
		select {
		case res = <-release:
			if susp == 0 {
				avail++
				free[res] = true
				fmt.Printf("[Server] Released resource: %d  \n", res)
			} else {
				blocked[i] = false
				susp--
				resource[i] <- res
				fmt.Printf("[Server] Woke process %d - Allocating resource %d.\n", i, res)
			}

		case p = <-when(avail > 0, request):
			for i = 0; i < MaxRes && !free[i]; i++ {
			}
			free[i] = false
			avail--
			resource[p] <- i
			fmt.Printf("[Server] Allocated resource %d to client %d.\n", i, p)

		case <-terminate:
			fmt.Println("END.")
			done <- 1
			return

		}
	}
}

func main() {
	var cli int
	fmt.Printf("\n How many clients? (max %d)? ", MaxProc)
	fmt.Scanf("%d", &cli)
	fmt.Println("Clients:", cli)

	for i := 0; i < cli; i++ {
		resource[i] = make(chan int, MaxBuff)
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
