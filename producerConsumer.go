package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Producer generates random integers and sends them to the shared channel.
func Producer(id int, dataChannel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		num := rand.Intn(100) // Generate a random number
		fmt.Printf("Producer %d: Produced %d\n", id, num)
		dataChannel <- num // Send the number to the channel
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	}
}

// Consumer reads integers from the shared channel and processes them (e.g., squares them).
func Consumer(id int, dataChannel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range dataChannel {
		fmt.Printf("Consumer %d: Consumed %d, Processed: %d\n", id, num, num*num)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dataChannel := make(chan int, 10) // Shared channel with a buffer size of 10
	var wg sync.WaitGroup

	// Start producers
	numProducers := 2
	for i := 1; i <= numProducers; i++ {
		wg.Add(1)
		go Producer(i, dataChannel, &wg)
	}

	// Start consumers
	numConsumers := 3
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		go Consumer(i, dataChannel, &wg)
	}

	// Wait for all producers to finish
	wg.Wait()

	// Close the data channel to signal consumers that no more data is coming
	close(dataChannel)

	// Give consumers time to finish processing remaining data
	time.Sleep(time.Second * 2)

	fmt.Println("All producers and consumers have finished.")
}
