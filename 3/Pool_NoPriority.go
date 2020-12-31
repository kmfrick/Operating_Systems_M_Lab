// resourcepool_prio project main.go

package main

import (
	"fmt"
	"time"
)

const MAXPROC = 10
const MAXRES = 3
const MAXBUFF = 20
const USAGE_TIME = 2

var request = make(chan int, MAXBUFF)
var release = make(chan int, MAXBUFF)
var resource [MAXPROC]chan int
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
	time.Sleep(time.Second * USAGE_TIME)
	release <- r
	done <- i
}

func server() {

	var avail int = MAXRES
	var res, p, i int
	var free [MAXRES]bool
	var susp = 0
	var blocked [MAXPROC]bool

	// init
	for i := 0; i < MAXRES; i++ {
		free[i] = true
	}
	for i := 0; i < MAXPROC; i++ {
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
			for i = 0; i < MAXRES && !free[i]; i++ {
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
	fmt.Printf("\n How many clients? (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("Clients:", cli)

	for i := 0; i < cli; i++ {
		resource[i] = make(chan int, MAXBUFF)
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
