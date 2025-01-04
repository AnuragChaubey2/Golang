# Golang


Problem 1: *Distributed Task Processing System with Timeout and Inter-Server Communication*

Problem: Design and implement a distributed task processing system where multiple servers (goroutines) process tasks concurrently. Each task has a processing deadline, and tasks that exceed their deadlines must be canceled. Servers communicate with each other via channels to simulate distributed communication, enabling task delegation, request handling, and result transmission. Use context.WithTimeout to enforce deadlines and ensure proper cancellation of tasks that exceed their processing time.
[ distributedSystem.go ]

Concepts:

Task scheduling and processing with deadlines
Goroutines and channels for concurrency and communication
context.WithTimeout for handling task timeouts
Simulated inter-server communication
