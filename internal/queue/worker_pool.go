package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dev-dhanushkumar/go-video-processor/internal/utils"
)

// WorkerPool manages a pool of workers for processing jobs
type WorkerPool struct {
	workerCount int
	queue       *JobQueue
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	logger      *utils.Logger
	activeJobs  sync.Map
	mu          sync.RWMutex
	stopped     bool
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int, queueSize int, logger *utils.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workerCount: workerCount,
		queue:       NewJobQueue(queueSize),
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
		stopped:     false,
	}
}

// Start starts all workers in the pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	if wp.stopped {
		wp.mu.Unlock()
		wp.logger.Warn("Worker pool already stopped, cannot start")
		return
	}
	wp.mu.Unlock()

	wp.logger.Info(fmt.Sprintf("Starting worker pool with %d workers", wp.workerCount))

	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the main worker goroutine
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	wp.logger.Info(fmt.Sprintf("Worker %d started", id))

	for {
		select {
		case <-wp.ctx.Done():
			wp.logger.Info(fmt.Sprintf("Worker %d stopping", id))
			return
		case job := <-wp.queue.Jobs():
			if job == nil {
				continue
			}

			wp.processJob(id, job)
		}
	}
}

// processJob processes a single job
func (wp *WorkerPool) processJob(workerID int, job *Job) {
	startTime := time.Now()

	wp.logger.Info(fmt.Sprintf("Worker %d processing job %s (operation: %s)", workerID, job.ID, job.Operation))

	// Track active job
	wp.activeJobs.Store(job.ID, job)
	defer wp.activeJobs.Delete(job.ID)

	// Create job context with timeout
	jobCtx, cancel := context.WithTimeout(wp.ctx, 30*time.Minute)
	defer cancel()

	// Execute job handler
	err := job.Handler(jobCtx)

	duration := time.Since(startTime)

	if err != nil {
		wp.logger.Error(fmt.Sprintf("Worker %d failed to process job %s: %v (duration: %s)",
			workerID, job.ID, err, duration))
	} else {
		wp.logger.Info(fmt.Sprintf("Worker %d completed job %s (duration: %s)",
			workerID, job.ID, duration))
	}
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job *Job) error {
	wp.mu.RLock()
	stopped := wp.stopped
	wp.mu.RUnlock()

	if stopped {
		return fmt.Errorf("worker pool is stopped")
	}

	return wp.queue.Enqueue(job)
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if wp.stopped {
		wp.mu.Unlock()
		return
	}
	wp.stopped = true
	wp.mu.Unlock()

	wp.logger.Info("Stopping worker pool...")

	// Cancel context to signal workers to stop
	wp.cancel()

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close the queue
	wp.queue.Close()

	wp.logger.Info("Worker pool stopped")
}

// ActiveJobs returns the number of currently active jobs
func (wp *WorkerPool) ActiveJobs() int {
	count := 0
	wp.activeJobs.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// QueueDepth returns the number of jobs waiting in the queue
func (wp *WorkerPool) QueueDepth() int {
	return wp.queue.Pending()
}

// WorkerCount returns the number of workers in the pool
func (wp *WorkerPool) WorkerCount() int {
	return wp.workerCount
}
