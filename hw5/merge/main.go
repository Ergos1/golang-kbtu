package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generator(from, to int) <-chan int {
	ch := make(chan int)

	go func() {
		for i := from; i <= to; i++ {
			ch <- i
			time.Sleep(time.Second * time.Duration(rand.Intn(3)))
		}
		close(ch)
	}()

	return ch
}

func merge(cs ...<-chan int) <-chan int {
	ch := make(chan int)
	wg := new(sync.WaitGroup)

	for _, c := range cs {
		wg.Add(1)
		localC := c
		go func() {
			defer wg.Done()

			for in := range localC {
				ch <- in
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func main() {
	var cs []<-chan int
	for i := 1; i <= 5; i++ {
		cs = append(cs, generator(i, i+1))
	}

	merged := merge(cs...)

	for res := range merged {
		fmt.Println(res)
	}

	fmt.Println("reading merge completed")
}
