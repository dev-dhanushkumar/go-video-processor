package queue

import (
	"context"
	"sync"
)

// Job represents a processing job
type Job struct {
	ID         string
	VideoID    string
	Operation  string
	Parameters map[string]interface{}
	Handler    func(context.Context) error
}

// JobQueue manages the queue of processing jobs
type JobQueue struct {
	jobs    chan *Job
	mu      sync.RWMutex
	size    int
	pending int
}

// NewJobQueue creates a new job queue with specified size
func NewJobQueue(size int) *JobQueue {
	return &JobQueue{
		jobs: make(chan *Job, size),
		size: size,
	}
}

// Enqueue adds a job to the queue
func (jq *JobQueue) Enqueue(job *Job) error {
	jq.mu.Lock()
	jq.pending++
	jq.mu.Unlock()

	select {
	case jq.jobs <- job:
		return nil
	default:
		jq.mu.Lock()
		jq.pending--
		jq.mu.Unlock()
		return ErrQueueFull
	}
}

// Dequeue removes and returns a job from the queue
func (jq *JobQueue) Dequeue() *Job {
	select {
	case job := <-jq.jobs:
		jq.mu.Lock()
		jq.pending--
		jq.mu.Unlock()
		return job
	default:
		return nil
	}
}

// Jobs returns the channel for receiving jobs
func (jq *JobQueue) Jobs() <-chan *Job {
	return jq.jobs
}

// Size returns the maximum size of the queue
func (jq *JobQueue) Size() int {
	return jq.size
}

// Pending returns the number of pending jobs
func (jq *JobQueue) Pending() int {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	return jq.pending
}

// Close closes the job queue
func (jq *JobQueue) Close() {
	close(jq.jobs)
}

var (
	ErrQueueFull = &QueueError{"queue is full"}
)

type QueueError struct {
	message string
}

func (e *QueueError) Error() string {
	return e.message
}
