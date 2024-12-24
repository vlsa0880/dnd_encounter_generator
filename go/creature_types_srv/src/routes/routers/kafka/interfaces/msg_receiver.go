package interfaces

type MsgReceiver interface {
	SetupMsgHandler(handler MsgHandler) error
}
