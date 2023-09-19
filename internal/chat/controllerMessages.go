package chat

const (
	CONTROLLER_REQUEST_JOIN uint8 = iota
	CONTROLLER_RESPONSE_JOIN
	// CONTROLLER_REQUEST_LEAVE
	// CONTROLLER_RESPONSE_LEAVE
)

const (
	RESPONSE_OK    uint8 = 1
	RESPONSE_ERROR uint8 = 2
)

func UnWrapControllerMessageRaw(bytes *[]byte) (ControllerMessage, error) {
	var raw ControllerMessageRaw
	raw.Type = (*bytes)[0]
	raw.Data = (*bytes)[1:]
	return raw, nil
}

type ControllerMessage interface {
	Wrap() []byte
}

type ControllerMessageRaw struct {
	Type uint8
	Data []byte
}

func (c ControllerMessageRaw) Wrap() []byte {
	data := make([]byte, 0, 1+len(c.Data))
	data = append(data, c.Type)
	data = append(data, c.Data...)

	return data
}

type ControllerRequestJoin struct {
	Username string
}

func (c ControllerRequestJoin) Wrap() []byte {
	raw := ControllerMessageRaw{
		Type: CONTROLLER_REQUEST_JOIN,
		Data: []byte(c.Username),
	}
	return raw.Wrap()
}

type ControllerResponseJoin struct {
	ResponseStatus uint8
	Message        *string
}

func (c ControllerResponseJoin) Wrap() []byte {
	raw := ControllerMessageRaw{
		Type: CONTROLLER_RESPONSE_JOIN,
		Data: []byte{c.ResponseStatus},
	}
	if c.Message != nil {
		raw.Data = append(raw.Data, []byte(*c.Message)...)
	}
	return raw.Wrap()
}
