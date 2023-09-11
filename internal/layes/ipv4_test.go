package layes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIpv4Layer(t *testing.T) {
	ipva := NewIpv4Layer()

	assert.Equal(t, uint8(4), ipva.Version, "Version should be 4")
	assert.Equal(t, uint8(5), ipva.IHL, "IHL should be 5")
}

func TestUnWrapIpv4(t *testing.T) {
	testsCases := []struct {
		name     string
		input    []byte
		expected Ipv4Layer
	}{
		{
			name:     "Empty data",
			input:    make([]byte, 20),
			expected: Ipv4Layer{Data: []byte{}},
		},
		{
			name:  "Valid TCP with no value",
			input: []byte{0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x06, 0x00, 0x00, 0x67, 0x83, 0x2b, 0x1f, 0xc1, 0x2a, 0x2a, 0x2a},
			expected: Ipv4Layer{
				Version:  4,
				IHL:      5,
				TOS:      0,
				Length:   0,
				ID:       0,
				Flags:    0,
				Fragment: 0,
				TTL:      64,
				Protocol: 0x06,
				Checksum: 0,
				Origem:   [4]byte{103, 131, 43, 31},
				Destino:  [4]byte{193, 42, 42, 42},
				Data:     []byte{},
			},
		},
	}

	for _, test := range testsCases {
		actual := UnWrapIpv4(&test.input)
		assert.Equal(t, test.expected, actual, "they should be equal")
	}
}

func TestIPV4ToBytes(t *testing.T) {

}

func TestCalculateChecksum(t *testing.T) {

}
