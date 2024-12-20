// A worker pool implementation to process tasks concurrently, where each task is the square of a number.

package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work to be processed by the workers.
type Task struct {
	Number int
}

// Worker represents a worker that processes tasks.
type Worker struct {
	ID      int
	Tasks   chan Task
	Quit    chan bool
	Results chan int
	wg      *sync.WaitGroup
}

// NewWorkerInstance creates and returns a new Worker.
func NewWorkerInstance(id int, tasks chan Task, results chan int, wg *sync.WaitGroup) Worker {
	return Worker{
		ID:      id,
		Tasks:   tasks,
		Quit:    make(chan bool),
		Results: results,
		wg:      wg,
	}
}

// StartWorker begins the worker's process to fetch tasks and process them.
func (w Worker) StartWorker() {
	go func() {
		defer w.wg.Done()
		for {
			select {
			case task, ok := <-w.Tasks:
				if !ok {
					return
				}
				result := task.Number * task.Number
				w.Results <- result
			case <-w.Quit:
				return
			}
		}
	}()
}

// Pool represents a pool of workers managing tasks and results.
type Pool struct {
	Workers    []Worker
	Tasks      chan Task
	Results    chan int
	NumWorkers int
}

// NewPoolInstance creates and returns a new WorkerPool.
func NewPoolInstance(numWorkers, numTasks int) *Pool {
	tasks := make(chan Task, numTasks)
	results := make(chan int, numTasks)
	return &Pool{
		Workers:    make([]Worker, numWorkers),
		Tasks:      tasks,
		Results:    results,
		NumWorkers: numWorkers,
	}
}

// StartPool initializes the workers, assigns tasks, and starts the task processing.
func (p *Pool) StartPool() {
	var wg sync.WaitGroup
	for i := 0; i < p.NumWorkers; i++ {
		p.Workers[i] = NewWorkerInstance(i+1, p.Tasks, p.Results, &wg)
		wg.Add(1)
		p.Workers[i].StartWorker()
	}

	// Add tasks to the tasks channel.
	for i := 1; i <= cap(p.Tasks); i++ {
		p.Tasks <- Task{Number: i}
	}

	close(p.Tasks)

	wg.Wait()

	close(p.Results)
}

// DisplayResults prints the results
func (p *Pool) DisplayResults() {
	for result := range p.Results {
		fmt.Println("Result:", result)
	}
}

// StopPool signals all workers to stop.
func (p *Pool) StopPool() {
	for _, worker := range p.Workers {
		worker.Quit <- true
	}
}

func main() {
	numWorkers := 3
	numTasks := 15

	// Create a pool of workers and start processing tasks.
	pool := NewPoolInstance(numWorkers, numTasks)
	pool.StartPool()

	// Display the processed results.
	pool.DisplayResults()

	// Gracefully stop all workers.
	pool.StopPool()

	// Allow time for workers to shut down before program exits.
	time.Sleep(time.Second)
}
