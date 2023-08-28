package layes

import "fmt"

const (
	IPV4_PROTOCOL_TCP = 0x06
	IPV4_PROTOCOL_UDP = 0x11
)

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
		Protocol: 0, // 0x06 TCP, 0x11 UDP
		Checksum: 0,
		Origem:   [4]byte{0, 0, 0, 0},
		Destino:  [4]byte{0, 0, 0, 0},
		Data:     []byte{},
	}
}

func UnWrapIpv4(bytes *[]byte) Ipv4Layer {
	var raw Ipv4Layer
	raw.Version = (*bytes)[0] >> 4
	raw.IHL = (*bytes)[0] & 0x0F
	raw.TOS = (*bytes)[1]
	raw.Length = uint16((*bytes)[2])<<8 | uint16((*bytes)[3])
	raw.ID = uint16((*bytes)[4])<<8 | uint16((*bytes)[5])
	raw.Flags = (*bytes)[6] >> 5
	raw.Fragment = uint16((*bytes)[6]&0x1F)<<8 | uint16((*bytes)[7])
	raw.TTL = (*bytes)[8]
	raw.Protocol = (*bytes)[9]
	raw.Checksum = uint16((*bytes)[10])<<8 | uint16((*bytes)[11])
	copy(raw.Origem[:], (*bytes)[12:16])
	copy(raw.Destino[:], (*bytes)[16:20])
	raw.Data = (*bytes)[20:]
	return raw
}

func UnWrapIpv4FromEthernet(eth EthernetLayer) Ipv4Layer {
	return UnWrapIpv4(&eth.Data)
}

func (i Ipv4Layer) String() string {
	str := "IPv4 Layer\n"
	str += fmt.Sprintf("Versão: %d \n", i.Version)
	str += fmt.Sprintf("IHL: %d \n", i.IHL)
	str += fmt.Sprintf("TOS: %d \n", i.TOS)
	str += fmt.Sprintf("Length: %d \n", i.Length)
	str += fmt.Sprintf("ID: %d \n", i.ID)
	str += fmt.Sprintf("Flags: %d \n", i.Flags)
	str += fmt.Sprintf("Fragment: %d \n", i.Fragment)
	str += fmt.Sprintf("TTL: %d \n", i.TTL)
	str += fmt.Sprintf("Protocol: %d \n", i.Protocol)
	str += fmt.Sprintf("Checksum: %d \n", i.Checksum)
	str += fmt.Sprintf("Origem: %d.%d.%d.%d \n", i.Origem[0], i.Origem[1], i.Origem[2], i.Origem[3])
	str += fmt.Sprintf("Destino: %d.%d.%d.%d \n", i.Destino[0], i.Destino[1], i.Destino[2], i.Destino[3])
	str += fmt.Sprintf("Data: %d Bytes \n\n", len(i.Data))

	return str
}
