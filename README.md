# Golang Problems


Problem 1: **Distributed Task Processing System with Timeout and Inter-Server Communication**

Problem: Design and implement a distributed task processing system where multiple servers (goroutines) process tasks concurrently. Each task has a processing deadline, and tasks that exceed their deadlines must be canceled. Servers communicate with each other via channels to simulate distributed communication, enabling task delegation, request handling, and result transmission. Use context.WithTimeout to enforce deadlines and ensure proper cancellation of tasks that exceed their processing time.
[ distributedSystem.go ]

Concepts:

- Task scheduling and processing with deadlines
- Goroutines and channels for concurrency and communication
- context.WithTimeout for handling task timeouts
- Simulated inter-server communication


Problem 2: **Producer-Consumer Problem with Channels**

Problem: Implement the producer-consumer problem using goroutines and channels. Have one or more producer goroutines that generate data (e.g., integers) and put them into a shared channel, and one or more consumer goroutines that take data from the channel and process it (e.g., print the integers squared). Ensure proper synchronization using channels and avoid race conditions.
[ producerConsumer.go ]

Concepts:

- Producer-consumer problem
- Goroutines
- Channels
- Synchronization


Problem 3: **Rate Limiter Using Goroutines and Mutex**

Problem: Implement a simple rate limiter that allows only n requests per second. Use a goroutine to periodically reset the rate limiter and a mutex to synchronize access to the counter that tracks the number of requests in the current second.
[ rateLimiter.go ]

Concepts:

- Goroutines
- sync.Mutex
- Rate limiting
- Timing

Problem 4: **Implement a Worker Pool for Task Processing**

Problem: Create a simple worker pool that processes a set of tasks concurrently. The task could be something simple like squaring integers. Each worker in the pool will pick tasks from a shared queue and process them. Use channels for communication between the main routine and workers, and ensure that the pool shuts down gracefully once all tasks are completed.

Concepts:

- Worker pool
- Goroutines
- Channels
- Task scheduling and synchronization


Problem 4 : **Graceful Shutdown of a Web Server**

Problem: Write a Go web server that allows graceful shutdown. When the server receives a shutdown signal (e.g., SIGINT), it should stop accepting new connections and finish processing any ongoing requests before shutting down. Use a context.Context to propagate cancellation signals to the ongoing requests.
[ gracefulShutDown.go ]

Concepts:
- Graceful shutdown
- context.WithCancel
- Handling OS signals
