package interfaces

type Consumer interface {
	Run() error
	Stop() error
}
