package infra

import (
	"context"
	"errors"
	"receiver/internal/domain/entities"
	"sync"
	"time"
)

type Pool struct {
	queue          chan entities.KafkaTask
	workers        int
	wg             sync.WaitGroup
	stop           chan struct{}
	infra          *Infra
	enqueueTimeout time.Duration
}

func NewPool(infra *Infra, workers int, queueSize int, timeout time.Duration) *Pool {
	return &Pool{
		queue:          make(chan entities.KafkaTask, queueSize),
		workers:        workers,
		wg:             sync.WaitGroup{},
		stop:           make(chan struct{}),
		infra:          infra,
		enqueueTimeout: timeout,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.queue {
				err := p.infra.Producer.SendMessage("test", job.Msg)
				if err != nil {
					job.Reply <- entities.KafkaResult{job.ID, false, err}
				} else {
					job.Reply <- entities.KafkaResult{job.ID, true, nil}
				}
			}
		}()
	}
}

func (p *Pool) Enqueue(msg entities.KafkaTask) error {
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
