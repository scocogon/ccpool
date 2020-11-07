package ccpool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrPoolStopped = errors.New("POOL was stopped")

type FnRACall func(context.Context, interface{}) interface{}

func NewResultArgPool(ctx context.Context, size int, fn FnRACall) ResultArgPool {
	p := newPool(ctx, size, fn)
	p.hasarg = true
	p.hasresult = true

	if p.opts.wm != nil {
		p.wm = p.opts.wm
	} else {
		p.wm = newWM(p.ctx, int(size))
	}

	return p
}

func (p *pool) Invoke(arg interface{}) (err error, result interface{}) {
	err, w := p.wm.GetWorker()
	if err != nil {
		return err, nil
	}

	return nil, w.exec(p.ctx, arg)

	// w.submit(p.ctx, arg)
	// return nil, w.result()
}

func (p *pool) stoprapool() {
	p.wm.Stop()
	p.wg.Done()
}

func (p *pool) startrapool() {
	n := int(p.cap)
	p.wg.Add(1)

	for i := 0; i < n; i++ {
		w := newWorker(p, p.wm)
		p.wm.AddWorker(w)

		if !p.hasresult {
			go w.run()
		}
	}

	select {
	case <-p.ctx.Done():
		println("call stoprapool")
		p.stoprapool()
	}
}

type pool struct {
	ctx    context.Context
	cancel func()

	wg sync.WaitGroup

	opts *Options

	size    int32
	cap     int32
	running int32

	hasarg    bool
	hasresult bool

	fn FnRACall
	wm WorkerManager
}

func (p *pool) Serve() {
	if !atomic.CompareAndSwapInt32(&p.running, 0, 1) {
		return
	}

	p.running = 1

	switch {
	case p.hasarg && p.hasresult:
		p.startrapool()

	case p.hasarg:
	default:
		p.startPool()
	}
}

func (p *pool) Wait()             { p.wg.Wait() }
func (p *pool) Stop()             { p.cancel() }
func (p *pool) addRunner(n int32) { atomic.AddInt32(&p.size, n) }

func newPool(ctx context.Context, cap int, fn FnRACall) *pool {
	if ctx == nil {
		ctx = context.Background()
	}

	var cancel func()
	ctx, cancel = context.WithCancel(ctx)

	var opt Options
	opt = *opts

	return &pool{
		ctx:    ctx,
		cancel: cancel,

		opts: &opt,

		size: 0,
		cap:  int32(cap),

		fn: fn,
	}
}
