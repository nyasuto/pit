package objects

import (
	"testing"

	"github.com/nyasuto/pit/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	data := []byte("Hello, World\n")
	obj := New(ObjectTypeBlob, data)

	assert.Equal(t, ObjectTypeBlob, obj.Type)
	assert.Equal(t, "blob 13\x00Hello, World\n", obj.Data)
	assert.Equal(t, hash.SHA1("3fa0d4b98289a95a7cd3a45c9545e622718f8d2b"), obj.Hash)
}
