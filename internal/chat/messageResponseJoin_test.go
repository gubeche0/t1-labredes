package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapAndUnWrapMessageResponseJoin(t *testing.T) {

	msg := MessageResponseJoin{
		UserName:  "UserName_qualquer",
		Succeeded: true,
	}

	bytes := msg.Wrap()

	msgUnwrap, err := UnWrapMessageResponseJoin(&bytes)

	assert.Nil(t, err)
	assert.Equal(t, msg.UserName, msgUnwrap.UserName)
	assert.Equal(t, msg.Succeeded, msgUnwrap.Succeeded)
	assert.Equal(t, uint32(0x18), msgUnwrap.MessageLen)
}
