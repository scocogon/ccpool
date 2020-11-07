package ccpool

type _pool interface {
	Serve()

	Wait()
	Stop()
}

type Pool interface {
	_pool
}

type ArgPool interface {
	_pool
}

type ResultArgPool interface {
	_pool

	Invoke(arg interface{}) (err error, result interface{})
}
