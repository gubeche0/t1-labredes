package chat

const (
	MESSAGE_TYPE_TEXT uint8 = iota
	MESSAGE_TYPE_FILE
	MESSAGE_TYPE_CONTROLLE
)

type Protocol interface {
	ListenMessage(port int) (<-chan MessageInterface, error)
	ListenControllerMessage(port int) (<-chan ControllerMessage, error)
	GetProtocolName() string
}
