package layes

import "fmt"

type UdpLayer struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
	Data            []byte
}

func NewUdpLayer() UdpLayer {
	return UdpLayer{
		SourcePort:      0,
		DestinationPort: 0,
		Length:          0,
		Checksum:        0,
		Data:            []byte{},
	}
}

func UnWrapUdp(bytes *[]byte) (*UdpLayer, error) {
	var raw UdpLayer
	raw.SourcePort = uint16((*bytes)[0])<<8 | uint16((*bytes)[1])
	raw.DestinationPort = uint16((*bytes)[2])<<8 | uint16((*bytes)[3])
	raw.Length = uint16((*bytes)[4])<<8 | uint16((*bytes)[5])
	raw.Checksum = uint16((*bytes)[6])<<8 | uint16((*bytes)[7])
	// copy(raw.Data[:], (*bytes)[8:])
	raw.Data = (*bytes)[8:]

	return &raw, nil
}

func UnWrapUdpFromIpv4(ipv4 Ipv4Layer) (*UdpLayer, error) {
	if ipv4.Protocol != IPV4_PROTOCOL_UDP {
		return nil, fmt.Errorf("Protocolo inválido")
	}

	return UnWrapUdp(&ipv4.Data)
}

func (u UdpLayer) String() string {
	str := "UDP Layer\n"
	str += fmt.Sprintf("Source Port: %d\n", u.SourcePort)
	str += fmt.Sprintf("Destination Port: %d\n", u.DestinationPort)
	str += fmt.Sprintf("Length: %d\n", u.Length)
	str += fmt.Sprintf("Checksum: %d. (Valid: %t)\n", u.Checksum, u.ChecksumIsValid())
	str += fmt.Sprintf("Data: %d Bytes \n", len(u.Data))

	return str
}

// @TODO: Implement
func (u UdpLayer) ChecksumIsValid() bool {
	return false
}

// @TODO: Implement
func (u *UdpLayer) CalculateChecksum() {
	u.Checksum = 0
}

func (u *UdpLayer) CalculateLength() {
	u.Length = uint16(len(u.Data) + 8) // 8 = UDP Header Size
}

func (u *UdpLayer) Prepare() {
	u.CalculateLength()
	u.CalculateChecksum()
}

func (u UdpLayer) ToBytes() []byte {
	var bytes []byte

	bytes = append(bytes, byte(u.SourcePort>>8))
	bytes = append(bytes, byte(u.SourcePort))

	bytes = append(bytes, byte(u.DestinationPort>>8))
	bytes = append(bytes, byte(u.DestinationPort))

	bytes = append(bytes, byte(u.Length>>8))
	bytes = append(bytes, byte(u.Length))

	bytes = append(bytes, byte(u.Checksum>>8))
	bytes = append(bytes, byte(u.Checksum))

	bytes = append(bytes, u.Data[:]...)

	return bytes
}
