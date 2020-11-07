package ccpool

import (
	"context"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fn := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				t.Log("Done")
				return

			default:
				time.Sleep(1 * time.Second)
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	p := NewPool(ctx, 10, fn)
	p.Serve()
}
