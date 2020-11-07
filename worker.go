package ccpool

import (
	"context"
)

type Message struct {
	ctx  context.Context
	body interface{}
}

type Worker struct {
	p  *pool
	wm WorkerManager

	hasresult bool

	st  chan struct{}
	msg chan *Message
}

func newWorker(p *pool, wm WorkerManager) *Worker {
	w := &Worker{
		p:         p,
		wm:        wm,
		hasresult: p.hasresult,
	}

	if !w.hasresult {
		w.st = make(chan struct{})
		w.msg = make(chan *Message)
	}

	return w
}

func (w *Worker) run() {
	for {
		select {
		case <-w.st:
			return

		case param := <-w.msg:
			w.p.fn(param.ctx, param.body)
			w.wm.AddWorker(w)
		}
	}
}

func (w *Worker) submit(ctx context.Context, arg interface{}) {
	w.msg <- &Message{ctx: ctx, body: arg}
}

func (w *Worker) exec(ctx context.Context, arg interface{}) interface{} {
	res := w.p.fn(ctx, arg)
	w.wm.AddWorker(w)
	return res
}

func (w *Worker) Stop() {
	if !w.hasresult {
		w.st <- struct{}{}
	}
}
