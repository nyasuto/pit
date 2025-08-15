package objects

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/nyasuto/pit/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func Test_NewBlob(t *testing.T) {
	data := []byte("Hello, World\n")
	obj := NewBlob(data)
	expectedHash, _ := hash.Parse("3fa0d4b98289a95a7cd3a45c9545e622718f8d2b")
	assert.Equal(t, ObjectTypeBlob, obj.Type)
	assert.Equal(t, []byte("blob 13\x00Hello, World\n"), obj.Data)
	assert.Equal(t, expectedHash, obj.Hash)
}

func Test_WriteBlob(t *testing.T) {
	data := []byte("Hello, World\n")
	obj := NewBlob(data)
	// Expected git object path: .git/objects/3f/a0d4...
	hex := obj.Hash.String()

	name, err := Write(obj)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(name, "/"+hex[:2]+"/"+hex[2:]), "unexpected object path: %s", name)
	defer func() {
		// Cleanup created files/directories (ignore errors)
		_ = os.Remove(name)
	}()
	// Read the stored object (zlib-compressed) and inflate
	raw, err := os.ReadFile(name)
	assert.NoError(t, err)

	zr, err := zlib.NewReader(bytes.NewReader(raw))
	assert.NoError(t, err)
	defer zr.Close()

	inflated, err := io.ReadAll(zr)
	assert.NoError(t, err)

	// Stored content should be the exact header+payload
	assert.Equal(t, obj.Data, inflated)

	assert.Equal(t, "3fa0d4b98289a95a7cd3a45c9545e622718f8d2b", hex)
	expectedHash, _ := hash.Parse("3fa0d4b98289a95a7cd3a45c9545e622718f8d2b")

	obj, err = Read(name)
	assert.NoError(t, err)
	assert.Equal(t, ObjectTypeBlob, obj.Type)
	assert.Equal(t, []byte("blob 13\x00Hello, World\n"), obj.Data)
	assert.Equal(t, expectedHash, obj.Hash)
}
