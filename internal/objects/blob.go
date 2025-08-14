package objects

import (
	"fmt"

	"github.com/nyasuto/pit/pkg/hash"
)

type ObjectType string

const (
	ObjectTypeBlob   ObjectType = "blob"
	ObjectTypeTree   ObjectType = "tree"
	ObjectTypeCommit ObjectType = "commit"
)

type object struct {
	Type ObjectType // Type of the object (e.g., "blob", "tree", "commit")
	Hash hash.SHA1  // SHA1 hash of the object
	Data string     // Raw data of the object
}

func New(t ObjectType, data []byte) object {
	size := len(data)
	header := fmt.Sprintf("%s %d\x00", t, size)
	content := header + string(data)
	h := hash.Hash([]byte(content))

	return object{
		Type: t,
		Hash: h,
		Data: content,
	}

}

/*
func Write(o object) ([]byte, error) {
	data := []byte("Hello, World!\n")

	err := os.WriteFile("example.txt", data, 0644)
	if err != nil {
		panic(err)
	}

}*/
