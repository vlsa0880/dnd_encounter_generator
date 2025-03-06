package interfaces

type Router interface {
	Run() error
	Stop() error
}
