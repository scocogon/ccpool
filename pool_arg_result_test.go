package ccpool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/scocogon/ccpool"
)

func TestRAPool(t *testing.T) {
	fn := func(ctx context.Context, arg interface{}) interface{} {
		time.Sleep(100 * time.Millisecond)
		return arg.(int32) * 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	_ = cancel
	p := ccpool.NewResultArgPool(ctx, 10, fn)
	go p.Serve()

	var wg sync.WaitGroup
	var sum, failed int32
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err, res := p.Invoke(int32(1))
			if err != nil {
				atomic.AddInt32(&failed, 1)
				return
			}

			atomic.AddInt32(&sum, res.(int32))
		}()
	}

	p.Wait()
	wg.Wait()
	t.Logf("succ = %d, failed = %d\n", sum/2, failed)
}
