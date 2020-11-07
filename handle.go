package ccpool

type _pool interface {
	Serve()
}

type Pool interface {
	_pool
}

type ArgPool interface {
	_pool
}

type ResultArgPool interface {
	_pool
}
