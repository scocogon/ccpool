package ccpool_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/scocogon/ccpool"
)

func TestPool(t *testing.T) {
	var res int32
	fn := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				atomic.AddInt32(&res, 1)
				return

			default:
				time.Sleep(1 * time.Second)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = cancel
	p := ccpool.NewPool(ctx, 10, fn)
	p.Serve()

	t.Log("Done 1: res =", res) // Done

	res = 900
	p = ccpool.NewPool(nil, 10, fn)
	go func() { time.Sleep(3 * time.Second); p.Stop() }()
	p.Serve()

	t.Log("Done 2: res =", res) // Done 2: res = 910
}
