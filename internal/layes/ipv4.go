package layes

type Ipv4Layer struct {
	Version  uint8
	IHL      uint8
	TOS      uint8
	Length   uint16
	ID       uint16
	Flags    uint8
	Fragment uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Origem   [4]byte
	Destino  [4]byte
	// Options  []byte
	Data []byte
}

func NewIpv4Layer() Ipv4Layer {
	return Ipv4Layer{
		Version:  4,
		IHL:      5,
		TOS:      0,
		Length:   0,
		ID:       0,
		Flags:    0,
		Fragment: 0,
		TTL:      64,
		Protocol: 0,
		Checksum: 0,
		Origem:   [4]byte{0, 0, 0, 0},
		Destino:  [4]byte{0, 0, 0, 0},
		Data:     []byte{},
	}
}
