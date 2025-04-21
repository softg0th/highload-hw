package infra

import (
	"context"
	"errors"
	"log"
	"storage/internal/entities"
	"storage/internal/repository"
	"sync"
	"time"
)

type Pool struct {
	queue          chan entities.Document
	workers        int
	wg             sync.WaitGroup
	stop           chan struct{}
	enqueueTimeout time.Duration
	repo           *repository.Repository
}

func NewPool(repo *repository.Repository, workers int, queueSize int, timeout time.Duration) *Pool {
	return &Pool{
		queue:          make(chan entities.Document, queueSize),
		workers:        workers,
		wg:             sync.WaitGroup{},
		stop:           make(chan struct{}),
		repo:           repo,
		enqueueTimeout: timeout,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.queue {
				err := p.repo.DB.InsertPostsMongoStream(context.Background(), job)
				if err != nil {
					log.Fatal("insert posts error:", err)
				} else {
					log.Println("insert posts success")
				}
			}
		}()
	}
}

func (p *Pool) Enqueue(msg entities.Document) error {
	select {
	case p.queue <- msg:
		return nil
	case <-p.stop:
		return errors.New("pool closed")
	case <-time.After(p.enqueueTimeout):
		return errors.New("timeout")
	}
}

func (p *Pool) Shutdown(ctx context.Context) error {
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
