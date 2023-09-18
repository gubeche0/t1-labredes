package layes

// @TODO: isso funciona?
func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b)-1; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if len(b)%2 != 0 {
		s += uint32(b[len(b)-1])
	}
	s = s>>16 + s&0xffff
	s += s >> 16
	return uint16(^s)
}
