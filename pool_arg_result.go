package ccpool

import (
	"context"
	"sync"
	"sync/atomic"
)

type FnCall func(context.Context)
type FnACall func(context.Context, interface{}) error
type FnRACall func(context.Context, interface{}) (error, interface{})

type pool struct {
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup

	running bool

	size int32
	cap  int32

	hasarg    bool
	hasresult bool

	fn FnRACall
}

func (p *pool) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *pool) Serve() {
	switch {
	case p.hasarg && p.hasresult:
	case p.hasarg:
	default:
		p.startPool()
	}
}

func (p *pool) addRunner(n int32) { atomic.AddInt32(&p.size, n) }

func newPool(ctx context.Context, cap int32, fn FnRACall) *pool {
	if ctx == nil {
		ctx = context.Background()
	}

	var cancel func()
	ctx, cancel = context.WithCancel(ctx)

	return &pool{
		ctx:    ctx,
		cancel: cancel,

		size: 0,
		cap:  cap,

		fn: fn,
	}
}

func NewResultArgPool(ctx context.Context, size int32, fn FnRACall) ResultArgPool {
	p := newPool(ctx, size, fn)
	p.hasarg = true
	p.hasresult = true

	return p
}
