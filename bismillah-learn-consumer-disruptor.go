package main

import (
	disruptor "bismillah-learn-consumer-disruptor/src/disruptor"
	kafkaconsumer "bismillah-learn-consumer-disruptor/src/kafkaConsumer"
	"fmt"
)

func destructor() {
	fmt.Println("Exiting App...")

}
func main() {
	defer destructor()
	defer disruptor.DeInit()
	defer kafkaconsumer.DeInit()
	fmt.Println("Bismillah")
	disruptor.Init()
	kafkaconsumer.Init()

}
