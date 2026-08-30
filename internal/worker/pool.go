package worker

import (
	"context"
	"sync"
)

type Job func(ctx context.Context)

type Pool struct {
	jobs   chan Job
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewPool(ctx context.Context, size int) *Pool {
	if size < 1 {
		size = 1
	}
	cctx, cancel := context.WithCancel(ctx)
	p := &Pool{
		jobs:   make(chan Job, size*4),
		cancel: cancel,
	}
	p.wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-cctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					job(cctx)
				}
			}
		}()
	}
	return p
}

func (p *Pool) Submit(job Job) bool {
	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}

func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
}
