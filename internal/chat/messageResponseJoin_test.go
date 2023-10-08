package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapAndUnWrapMessageJoinRequestResponse(t *testing.T) {

	msg := MessageJoinRequestResponse{
		UserName:  "UserName_qualquer",
		Succeeded: true,
	}

	bytes := msg.Wrap()

	msgUnwrap, err := UnWrapMessageJoinRequestResponse(&bytes)

	assert.Nil(t, err)
	assert.Equal(t, msg.UserName, msgUnwrap.UserName)
	assert.Equal(t, msg.Succeeded, msgUnwrap.Succeeded)
	assert.Equal(t, uint64(0x1c), msgUnwrap.MessageLen)
}
