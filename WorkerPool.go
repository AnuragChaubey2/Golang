package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work to be processed by the workers.
type Task struct {
	ID        int
	Work      func(ctx context.Context) (int, error)
	Completed chan struct{}      
}

// Worker represents a worker that processes tasks.
type Worker struct {
	ID      int
	Tasks   chan Task
	Quit    chan bool
	Results chan int
	wg      *sync.WaitGroup
}

// WorkerPool manages a pool of worker goroutines.
type WorkerPool struct {
	Workers    []Worker
	Tasks      chan Task
	Results    chan int
	NumWorkers int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
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
func (w Worker) StartWorker(ctx context.Context) {
	go func() {
		defer w.wg.Done()
		for {
			select {
			case task, ok := <-w.Tasks:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d: task %d canceled\n", w.ID, task.ID)
				default:
					result, err := task.Work(ctx)
					if err == nil {
						w.Results <- result
					} else {
						fmt.Printf("Worker %d: error processing task %d: %v\n", w.ID, task.ID, err)
					}
				}
				close(task.Completed)
			case <-w.Quit:
				fmt.Printf("Worker %d: shutting down\n", w.ID)
				return
			}
		}
	}()
}

// NewWorkerPool creates and returns a new WorkerPool.
func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		Workers:    make([]Worker, numWorkers),
		Tasks:      make(chan Task),
		Results:    make(chan int),
		NumWorkers: numWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// StartPool initializes the workers and begins task processing.
func (wp *WorkerPool) StartPool() {
	for i := 0; i < wp.NumWorkers; i++ {
		worker := NewWorkerInstance(i+1, wp.Tasks, wp.Results, &wp.wg)
		wp.Workers[i] = worker
		wp.wg.Add(1)
		worker.StartWorker(wp.ctx)
	}
}

// AddTask adds a new task to the pool.
func (wp *WorkerPool) AddTask(task Task) {
	wp.Tasks <- task
}

// StopPool signals all workers to stop and waits for them to finish.
func (wp *WorkerPool) StopPool() {
	close(wp.Tasks) 
	wp.cancel()  
	for _, worker := range wp.Workers {
		worker.Quit <- true
	}
	wp.wg.Wait()
	close(wp.Results) 
}

func (wp *WorkerPool) DisplayResults() {
	for result := range wp.Results {
		fmt.Println("Result:", result)
	}
}

func main() {
	numWorkers := 3

	pool := NewWorkerPool(numWorkers)
	pool.StartPool()

	for i := 1; i <= 10; i++ {
		taskID := i
		task := Task{
			ID: taskID,
			Work: func(ctx context.Context) (int, error) {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				case <-time.After(2 * time.Second): 
					fmt.Printf("Task %d completed\n", taskID)
					return taskID * taskID, nil
				}
			},
			Completed: make(chan struct{}),
		}
		pool.AddTask(task)
	}

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Cancelling remaining tasks...")
		pool.StopPool()
	}()

	pool.DisplayResults()

	fmt.Println("All workers shut down.")
}
