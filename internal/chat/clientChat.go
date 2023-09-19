package chat

import (
	"errors"
	"log"
	"time"
)

type Conn interface {
	Conection() error
	GetReadMessageChannel() <-chan MessageInterface
	GetReadControllerMessageChannel() <-chan ControllerMessage
	GetWriteMessageChannel() chan<- MessageInterface
	GetWriteControllerMessageChannel() chan<- ControllerMessage
	Close() error
}

type ClientChat struct {
	Conn Conn
	User string
	// Address     string
	// MessagePort int
	// CommandPort int

	readMessageChannel            <-chan MessageInterface
	readControllerMessageChannel  <-chan ControllerMessage
	writeMessageChannel           chan<- MessageInterface
	writeControllerMessageChannel chan<- ControllerMessage
}

// ReciveUpdate
// SendMessage
// SendFile
// Join
// Leave

func (c *ClientChat) Connect() error {

	if c.User == "" {
		return errors.New("User is empty")
	}

	c.readMessageChannel = c.Conn.GetReadMessageChannel()
	c.readControllerMessageChannel = c.Conn.GetReadControllerMessageChannel()
	c.writeMessageChannel = c.Conn.GetWriteMessageChannel()
	c.writeControllerMessageChannel = c.Conn.GetWriteControllerMessageChannel()

	c.writeControllerMessageChannel <- ControllerRequestJoin{
		Username: c.User,
	}

	select {
	case message := <-c.readControllerMessageChannel:
		switch message.(type) {
		case ControllerResponseJoin:
			if message.(ControllerResponseJoin).ResponseStatus != RESPONSE_OK {
				return errors.New("Error to join")
			}
		default:
			return errors.New("Error to join, invalid response")
		}

	case <-time.After(5 * time.Second):
		return errors.New("Error to join, timeout")
	}

	log.Printf("Connected to server with user %s", c.User)

	go c.handleMessage()
	go c.handleControllerMessage()

	return nil
}

func (c ClientChat) handleMessage() {
	for message := range c.readMessageChannel {
		log.Printf("Recive message %v", message)
	}
}

func (c ClientChat) handleControllerMessage() {

	for message := range c.readControllerMessageChannel {
		log.Printf("Recive controller message %v", message)
	}
}

func (c ClientChat) SendMessage(message MessageInterface) {

	c.writeMessageChannel <- message
}
