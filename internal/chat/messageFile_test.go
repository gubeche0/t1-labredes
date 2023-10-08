package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapAndUnWrapMessageFile(t *testing.T) {

	msg := MessageFile{
		Origin:   "Origin_qualquer",
		Target:   "Target_qualquer",
		Filename: "Filename_qualquer.txt",
		// Filesize: 0x1234567890,
		Filedata: []byte("Filedata_qualquer"),
	}

	bytes := msg.Wrap()

	msgUnwrap, err := UnWrapMessageFile(&bytes)

	assert.Nil(t, err)
	assert.Equal(t, msg.Origin, msgUnwrap.Origin)
	assert.Equal(t, msg.Target, msgUnwrap.Target)
	assert.Equal(t, msg.Filename, msgUnwrap.Filename)
	assert.Equal(t, uint64(len(msg.Filedata)), msgUnwrap.Filesize)
	assert.Equal(t, msg.Filedata, msgUnwrap.Filedata)
}

func TestUnWrapMessageFile(t *testing.T) {
	bytes := []byte{
		5,
		0, 0, 0, 0, 0, 0, 0, 70,
		65, 110, 111, 110, 121, 109, 111, 117, 115, 95, 52, 55, 54, 10,
		65, 76, 76, 10,
		82, 69, 65, 68, 77, 69, 46, 109, 100, 10,
		0, 0, 0, 0, 0, 0, 0, 24,
		35, 32, 114, 97, 119, 45, 115, 111, 99, 107, 101, 116, 45, 116, 49, 45, 108, 97, 98, 114, 101, 100, 101, 115,
	}

	msg, err := UnWrapMessageFile(&bytes)

	assert.Nil(t, err)

	assert.Equal(t, MESSAGE_TYPE_FILE, msg.GetType())
	assert.Equal(t, uint64(70), msg.MessageLen)
	assert.Equal(t, "Anonymous_476", msg.Origin)
	assert.Equal(t, MESSAGE_TARGET_ALL, msg.Target)
	assert.Equal(t, "README.md", msg.Filename)
	assert.Equal(t, uint64(24), msg.Filesize)
	assert.Equal(t, []byte("# raw-socket-t1-labredes"), msg.Filedata)

}
