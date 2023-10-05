package disruptor

import (
	"fmt"
)

func DeInit() {
	fmt.Println("This is DeInit of Disruptor")

}
func Init() {
	// defer DeInit()
	fmt.Println("This is Init of Disruptor")

}
