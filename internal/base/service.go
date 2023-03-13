package base

type Servicer interface {
	Start() error
	Stop() error
}
