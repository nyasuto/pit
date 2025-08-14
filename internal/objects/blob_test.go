package objects

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func Test_Write(t *testing.T) {
	data := []byte("Hello, World\n")
	obj := New(ObjectTypeBlob, data)
	// Expected git object path: .git/objects/3f/a0d4...
	hex := obj.Hash.String()
	dir := filepath.Join(".test-git", "objects", hex[:2])
	path := filepath.Join(dir, hex[2:])

	name, err := Write(obj)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(name, "/"+hex[:2]+"/"+hex[2:]), "unexpected object path: %s", name)
	defer func() {
		// Cleanup created files/directories (ignore errors)
		_ = os.Remove(name)
	}()
	// Read the stored object (zlib-compressed) and inflate
	raw, err := os.ReadFile(path)
	assert.NoError(t, err)

	zr, err := zlib.NewReader(bytes.NewReader(raw))
	assert.NoError(t, err)
	defer zr.Close()

	inflated, err := io.ReadAll(zr)
	assert.NoError(t, err)

	// Stored content should be the exact header+payload
	assert.Equal(t, obj.Data, string(inflated))

	assert.Equal(t, "3fa0d4b98289a95a7cd3a45c9545e622718f8d2b", hex)

	obj, err = Read(name)
	assert.NoError(t, err)
	assert.Equal(t, ObjectTypeBlob, obj.Type)
	assert.Equal(t, "blob 13\x00Hello, World\n", obj.Data)
	assert.Equal(t, hash.SHA1("3fa0d4b98289a95a7cd3a45c9545e622718f8d2b"), obj.Hash)
}
