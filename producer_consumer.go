package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	bufferSize  = 5  // 缓冲区容量，写满后生产者阻塞
	producerNum = 3  // 生产者数量
	consumerNum = 2  // 消费者数量
	itemsEach   = 10 // 每个生产者生产的数量
)

func main() {
	// 单缓冲区：用带容量的 channel 实现，满了 send 会阻塞，空了 recv 会阻塞
	buffer := make(chan int, bufferSize)

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// 启动多个生产者
	for p := 0; p < producerNum; p++ {
		producerWg.Add(1)
		go func(id int) {
			defer producerWg.Done()
			for i := 0; i < itemsEach; i++ {
				item := id*100 + i
				buffer <- item // 缓冲区满时在此阻塞
				fmt.Printf("生产者 %d 生产: %d\n", id, item)
			}
		}(p)
	}

	// 启动多个消费者
	for c := 0; c < consumerNum; c++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			for item := range buffer { // 缓冲区空时在此阻塞；channel 关闭后退出
				fmt.Printf("  消费者 %d 消费: %d\n", id, item)
				time.Sleep(50 * time.Millisecond) // 模拟消费耗时，更容易观察到阻塞
			}
		}(c)
	}

	// 等所有生产者结束后关闭缓冲区，通知消费者收尾
	go func() {
		producerWg.Wait()
		close(buffer)
	}()

	consumerWg.Wait()
	fmt.Println("全部完成")
}
