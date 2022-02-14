package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	TestStr := "The string for testing"
	var err error

	var TestEncode string
	var TestDecode string

	// encryption test
	TestEncode, err = Encode(TestStr)
	assert.Empty(t, err)

	// decrypt test
	TestDecode, err = Decode(TestEncode)
	assert.Empty(t, err)

	// Verify decryption result
	assert.Equal(t, string(TestDecode), TestStr)
}
