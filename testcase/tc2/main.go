package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 150; i++ {
		time.Sleep(time.Second)
		fmt.Printf("id:[2] - %d \n", (i + 1))
	}
}
