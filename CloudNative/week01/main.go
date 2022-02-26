package main

import (
	"fmt"
	"time"
)

func productor(queue chan int) {
	for i := 0; i < 20; i++ {
		fmt.Printf("生产消息: %d\n", i)
		queue <- i
		time.Sleep(time.Second)
	}
}

func consumer(queue chan int) {
	for {
		res := <-queue
		fmt.Printf("消费消息: %d\n", res)
		time.Sleep(time.Second * 2)
	}
}

func main() {
	queue := make(chan int, 10)
	go productor(queue)
	consumer(queue)
	// time.Sleep(time.Second * 20)
	// close(queue)
}
