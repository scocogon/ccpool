package ccpool

type Options struct {
	MaxCapacity uint32

	// 自定义 worker 管理器
	wm WorkerManager
}

var opts = &Options{}
