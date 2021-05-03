package main

import (
	"context"
	"log"
	"math/rand"

	// "sync"
	"time"
)

func main() {

	root := context.Background()
	ctx, cancel := context.WithCancel(root)
	ch := make(chan int)

	// log.Print(ctx.Deadline())
	// <-ctx.Done()
	// log.Print(ctx.Deadline())
	// log.Print(ctx.Err())
	// log.Print(ctx.Err() == context.DeadlineExceeded)

	for i := 0; i < 5; i++ {
		go func(ctx context.Context, index int, ch chan<- int) {
			wait := rand.Intn(10)
			log.Printf("%d goroutine wait %d-second", index, wait)
			select {
			case <-ctx.Done():
				log.Printf("%d canceled", index)
			case <-time.After(time.Second * time.Duration(wait)):
				ch <- index
				log.Printf("%d goroutine done with %d-seconds ", index, wait)
			}
		}(ctx, i, ch)
	}

	winner := <-ch
	cancel()
	log.Print(winner)
	<-time.After(time.Second)
	log.Print("done")

	// wg := sync.WaitGroup{}
	// wg.Add(1)

	// go func(ctx context.Context) {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Print("done")
	// 			wg.Done()
	// 			return
	// 		case <-time.After(time.Second):
	// 			log.Print("tick")
	// 		}
	// 	}
	// }(ctx)

	// wg.Wait()
	// log.Print("main done")

}
