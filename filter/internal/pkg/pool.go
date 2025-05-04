package pkg

import (
	"context"
	"errors"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"sync"
	"time"
)

type Pool[T any] struct {
	queue          chan T
	workers        int
	wg             sync.WaitGroup
	stop           chan struct{}
	enqueueTimeout time.Duration
	runningFunc    func(context.Context, T, interface{}) error
	metadata       interface{}
	logger         *logstash_logger.Logstash
}

func NewPool[T any](workers int, queueSize int, timeout time.Duration,
	runningFunc func(context.Context, T, interface{}) error, metadata interface{},
	logger *logstash_logger.Logstash) *Pool[T] {
	return &Pool[T]{
		queue:          make(chan T, queueSize),
		workers:        workers,
		wg:             sync.WaitGroup{},
		stop:           make(chan struct{}),
		enqueueTimeout: timeout,
		runningFunc:    runningFunc,
		metadata:       metadata,
		logger:         logger,
	}
}

func (p *Pool[T]) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.queue {
				err := p.runningFunc(context.Background(), job, p.metadata)
				if err != nil {
					p.logger.Error(map[string]interface{}{
						"message":           "Insert posts error",
						"error":             true,
						"error_description": err.Error(),
					})
				} else {
					p.logger.Info(map[string]interface{}{
						"message": "Insert posts success",
						"error":   false,
					})
				}
			}
		}()
	}
}

func (p *Pool[T]) Enqueue(msg T) error {
	select {
	case p.queue <- msg:
		return nil
	case <-p.stop:
		return errors.New("pool closed")
	case <-time.After(p.enqueueTimeout):
		return errors.New("timeout")
	}
}

func (p *Pool[T]) Shutdown(ctx context.Context) error {
	close(p.stop)
	close(p.queue)
	done := make(chan struct{})

	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
