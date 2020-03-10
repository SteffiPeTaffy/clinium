package clinium

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_matchesPrompt_removes_userInput_on_match(t *testing.T) {
	userInput := []UserInput{
		{
			Prompt: "Foo",
			Input:  "Bar",
		},
		{
			Prompt: "Cat",
			Input:  "Meow",
		},
	}

	_, input, err := indexOf("Foo", userInput)

	assert.Nil(t, err)
	assert.Equal(t, input.Input, "Bar")

}
