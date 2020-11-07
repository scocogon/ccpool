package ccpool

import (
	"context"
	"sync"
)

type WorkerManager interface {
	GetWorker() (error, *Worker)
	AddWorker(*Worker)

	Stop()
}

type wm struct {
	wg  sync.WaitGroup
	ctx context.Context
	ch  chan *Worker
}

func newWM(ctx context.Context, size int) WorkerManager {
	m := &wm{
		ctx: ctx,
		ch:  make(chan *Worker, size),
	}

	m.wg.Add(size)
	return m
}

func (m *wm) GetWorker() (error, *Worker) {
	select {
	case <-m.ctx.Done():
		return ErrPoolStopped, nil

	case w, ok := <-m.ch:
		if ok {
			return nil, w
		}

		return ErrPoolStopped, nil
	}
}

func (m *wm) AddWorker(w *Worker) {
	m.ch <- w
}

func (m *wm) Stop() {
	go func() {
		for w := range m.ch {
			w.Stop()
			m.wg.Done()
		}
	}()

	m.wg.Wait()
	close(m.ch)
}
