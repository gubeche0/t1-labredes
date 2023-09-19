package chat

type ServerChat struct {
	Protocol    Protocol
	Users       []string
	Address     string
	MessagePort int
	CommandPort int
	//
}

// Listen
// Handlers
// StartListenAndServer

func (c *ServerChat) StartListenAndServer() error {

	return nil
}

func (c ServerChat) handleMessageRaw(message MessageRaw) {}

func (c ServerChat) handleControllerMessageRaw(message ControllerMessageRaw) {}
