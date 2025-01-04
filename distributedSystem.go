package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work with a unique ID and processing time.
type Task struct {
	ID            int
	ProcessingTime time.Duration
}

// Server represents a worker server that processes tasks and communicates with other servers.
type Server struct {
	ID      int
	Tasks   chan Task       // Channel for receiving tasks
	Results chan string     // Channel for sending results
	Peers   []chan Task     // Channels to communicate with other servers
	wg      *sync.WaitGroup // WaitGroup for synchronization
}

// NewServer creates a new server instance.
func NewServer(id int, numPeers int, wg *sync.WaitGroup) *Server {
	return &Server{
		ID:      id,
		Tasks:   make(chan Task),
		Results: make(chan string),
		Peers:   make([]chan Task, numPeers),
		wg:      wg,
	}
}

// Start starts the server's task processing loop.
func (s *Server) Start() {
	go func() {
		defer s.wg.Done()
		for task := range s.Tasks {
			ctx, cancel := context.WithTimeout(context.Background(), task.ProcessingTime)
			defer cancel()

			result := make(chan string)
			go processTask(ctx, task, result)

			select {
			case res := <-result:
				s.Results <- fmt.Sprintf("Server %d: Completed Task %d -> %s", s.ID, task.ID, res)
			case <-ctx.Done():
				s.Results <- fmt.Sprintf("Server %d: Task %d timed out -> %s", s.ID, task.ID, ctx.Err())
			}
		}
	}()
}

// processTask simulates processing a task.
func processTask(ctx context.Context, task Task, result chan<- string) {
	select {
	case <-time.After(task.ProcessingTime): // Simulate task processing
		result <- fmt.Sprintf("Processed Task %d", task.ID)
	case <-ctx.Done():
		// Handle cancellation
		return
	}
}

// TaskGenerator generates tasks and assigns them to a random server.
func TaskGenerator(servers []*Server, numTasks int) {
	for i := 1; i <= numTasks; i++ {
		task := Task{
			ID:            i,
			ProcessingTime: time.Duration(rand.Intn(500)+100) * time.Millisecond,
		}
		// Assign task to a random server
		targetServer := servers[rand.Intn(len(servers))]
		fmt.Printf("TaskGenerator: Assigning Task %d to Server %d with processing time %v\n", task.ID, targetServer.ID, task.ProcessingTime)
		targetServer.Tasks <- task
	}
	for _, server := range servers {
		close(server.Tasks) // Close task channels to signal no more tasks
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numServers := 3
	numTasks := 10
	var wg sync.WaitGroup

	// Create servers
	servers := make([]*Server, numServers)
	for i := 0; i < numServers; i++ {
		server := NewServer(i+1, numServers, &wg)
		servers[i] = server
		wg.Add(1)
		server.Start()
	}

	// Set up inter-server communication (simulated by assigning peer channels)
	for i, server := range servers {
		for j, peer := range servers {
			if i != j {
				server.Peers[j] = peer.Tasks
			}
		}
	}

	// Generate and distribute tasks
	go TaskGenerator(servers, numTasks)

	// Collect results
	go func() {
		wg.Wait()
		for _, server := range servers {
			close(server.Results)
		}
	}()

	// Print results
	for _, server := range servers {
		go func(s *Server) {
			for result := range s.Results {
				fmt.Println(result)
			}
		}(server)
	}

	// Allow time for processing to complete
	time.Sleep(3 * time.Second)
	fmt.Println("All servers have completed their tasks.")
}