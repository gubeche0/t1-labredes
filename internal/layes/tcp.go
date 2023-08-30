package layes

import "fmt"

type TcpLayer struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	AckNumber       uint32
	DataOffset      uint8 // 4 bits
	Reserved        uint8 // 4 bits
	Flags           uint8
	WindowSize      uint16
	Checksum        uint16
	UrgentPointer   uint16
	// OptionalData    []byte
	Data []byte
}

func NewTcpLayer() *TcpLayer {
	return &TcpLayer{}
}

func UnWrapTcp(bytes *[]byte) (*TcpLayer, error) {
	var tcp TcpLayer
	tcp.SourcePort = uint16((*bytes)[0])<<8 | uint16((*bytes)[1])
	tcp.DestinationPort = uint16((*bytes)[2])<<8 | uint16((*bytes)[3])
	tcp.SequenceNumber = uint32((*bytes)[4])<<24 | uint32((*bytes)[5])<<16 | uint32((*bytes)[6])<<8 | uint32((*bytes)[7])
	tcp.AckNumber = uint32((*bytes)[8])<<24 | uint32((*bytes)[9])<<16 | uint32((*bytes)[10])<<8 | uint32((*bytes)[11])
	tcp.DataOffset = (*bytes)[12] >> 4
	tcp.Reserved = (*bytes)[12] & 0x0F
	tcp.Flags = (*bytes)[13]
	tcp.WindowSize = uint16((*bytes)[14])<<8 | uint16((*bytes)[15])
	tcp.Checksum = uint16((*bytes)[16])<<8 | uint16((*bytes)[17])
	tcp.UrgentPointer = uint16((*bytes)[18])<<8 | uint16((*bytes)[19])

	tcp.Data = (*bytes)[20:]

	return &tcp, nil
}

func UnWrapTcpFromIpv4(ipv4 Ipv4Layer) (*TcpLayer, error) {
	if ipv4.Protocol != IPV4_PROTOCOL_TCP {
		return nil, fmt.Errorf("Not a TCP packet")
	}
	return UnWrapTcp(&ipv4.Data)
}

func (tcp TcpLayer) String() string {
	str := "TCP Layer\n"
	str += fmt.Sprintf("Source Port: %d\n", tcp.SourcePort)
	str += fmt.Sprintf("Destination Port: %d\n", tcp.DestinationPort)
	str += fmt.Sprintf("Sequence Number: %d\n", tcp.SequenceNumber)
	str += fmt.Sprintf("Ack Number: %d\n", tcp.AckNumber)
	str += fmt.Sprintf("Data Offset: %d\n", tcp.DataOffset)
	// str += fmt.Sprintf("Reserved: %d\n", tcp.Reserved)
	str += fmt.Sprintf("Flags: %d\n", tcp.Flags)
	str += fmt.Sprintf("Window Size: %d\n", tcp.WindowSize)
	str += fmt.Sprintf("Checksum: %d. (Valid: %t)\n", tcp.Checksum, tcp.ChecksumIsValid())
	str += fmt.Sprintf("Urgent Pointer: %d\n", tcp.UrgentPointer)
	str += fmt.Sprintf("Data: %d Bytes \n", len(tcp.Data))

	return str
}

// @TODO: Implement
func (tcp TcpLayer) ChecksumIsValid() bool {
	return false
}

// @TODO: Implement
func (tcp *TcpLayer) CalculateChecksum() {
	tcp.Checksum = 0
}

// @TODO: Implement
func (tcp *TcpLayer) Prepare() {
	tcp.CalculateChecksum()
}

func (tcp TcpLayer) ToBytes() []byte {
	var bytes []byte

	bytes = append(bytes, byte(tcp.SourcePort>>8))
	bytes = append(bytes, byte(tcp.SourcePort))

	bytes = append(bytes, byte(tcp.DestinationPort>>8))
	bytes = append(bytes, byte(tcp.DestinationPort))

	bytes = append(bytes, byte(tcp.SequenceNumber>>24))
	bytes = append(bytes, byte(tcp.SequenceNumber>>16))
	bytes = append(bytes, byte(tcp.SequenceNumber>>8))
	bytes = append(bytes, byte(tcp.SequenceNumber))

	bytes = append(bytes, byte(tcp.AckNumber>>24))
	bytes = append(bytes, byte(tcp.AckNumber>>16))
	bytes = append(bytes, byte(tcp.AckNumber>>8))
	bytes = append(bytes, byte(tcp.AckNumber))

	bytes = append(bytes, byte(tcp.DataOffset<<4|tcp.Reserved))

	bytes = append(bytes, byte(tcp.Flags))

	bytes = append(bytes, byte(tcp.WindowSize>>8))
	bytes = append(bytes, byte(tcp.WindowSize))

	bytes = append(bytes, byte(tcp.Checksum>>8))
	bytes = append(bytes, byte(tcp.Checksum))

	bytes = append(bytes, byte(tcp.UrgentPointer>>8))
	bytes = append(bytes, byte(tcp.UrgentPointer))

	bytes = append(bytes, tcp.Data...)

	return bytes
}
