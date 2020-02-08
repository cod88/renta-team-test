package main

import (
	"fmt"
	"messenger"
	"time"
)

func main() {
	fmt.Println("Storage service")
	messenger.WaitQueryForNews()
	ticker := time.NewTicker(time.Second * 5)

	for _ = range ticker.C {
		fmt.Println("Waiting for queryes...")
	}
}
