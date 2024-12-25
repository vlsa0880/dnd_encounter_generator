package interfaces

type Consumer interface {
	MsgReceiver
	Run()
	Stop()
}
