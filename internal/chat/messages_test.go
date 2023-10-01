package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapAndUnWrap(t *testing.T) {

	msg := MessageText{
		Origin: "origin",
		Target: "target",
		Text:   "text qualquer",
	}

	bytes := msg.Wrap()

	msgUnwrap, err := UnWrapMessageText(&bytes)

	assert.Nil(t, err)
	assert.Equal(t, msg.Origin, msgUnwrap.Origin)
	assert.Equal(t, msg.Target, msgUnwrap.Target)
	assert.Equal(t, msg.Text, msgUnwrap.Text)
	// assert.Equal(t, msg, *msgUnwrap)
}
