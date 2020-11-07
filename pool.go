package ccpool

import (
	"context"
)

func NewPool(ctx context.Context, size int32, fn FnCall) Pool {
	p := newPool(ctx, size, nil)
	p.fn = func(ctx context.Context, _ interface{}) (error, interface{}) {
		fn(ctx)
		p.wg.Done()
		return nil, nil
	}

	return p
}

func (p *pool) startPool() {
	n := int(p.cap)
	p.wg.Add(n)

	for i := 0; i < n; i++ {
		go p.fn(p.ctx, nil)
	}

	p.wg.Wait()
}
