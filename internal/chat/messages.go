package chat

type MessageInterface interface {
	// UnWrap(bytes *[]byte) Message
	Wrap() []byte
	GetTarget() string
	GetOrigin() string
	GetType() uint8
}

func UnWrapMessageRaw(bytes *[]byte) (MessageInterface, error) {
	var raw MessageRaw
	raw.Type = (*bytes)[0]
	raw.Rawdata = (*bytes)[1:]
	return raw, nil
}

type MessageRaw struct {
	Type    uint8
	Target  string
	Origin  string
	Rawdata []byte
}

func (m MessageRaw) Wrap() []byte {
	return []byte{}
}

func (m MessageRaw) GetType() uint8 {
	return m.Type
}

func (m MessageRaw) GetTarget() string {
	return m.Target
}

func (m MessageRaw) GetOrigin() string {
	return m.Origin
}

type MessageText struct {
	MessageRaw
	Text string
}

type MessageFile struct {
	MessageRaw
	Filename string
	Filesize uint32
	Filedata []byte
}
