package objects

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nyasuto/pit/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func Test_NewTree(t *testing.T) {
	trees := []TreeEntry{
		{Name: "file1.txt", Hash: hash.SHA1{}, Mode: 0644},
		{Name: "dir1", Hash: hash.SHA1{}, Mode: 040000},
	}

	tree := NewTree()
	for _, entry := range trees {
		tree.AddEntry(entry)
	}
	obj := tree.toObject()

	assert.NotNil(t, obj)
	assert.Equal(t, ObjectTypeTree, obj.Type)

	//	if len(tree.Entries) != len(trees) {
	//		t.Errorf("Expected %d entries, got %d", len(trees), len(tree.Entries))
	//	}
}

func Test_SerializeTree(t *testing.T) {
	tree := NewTree()
	tree.AddEntry(TreeEntry{Name: "file1.txt", Hash: hash.SHA1{}, Mode: ModeFile})
tree.AddEntry(TreeEntry{Name: "dir1", Hash: hash.SHA1{}, Mode: ModeDir})

	obj := tree.serialize()
	assert.Equal(t, ObjectTypeTree, obj.Type)

	// Gitの正しい形式: mode name\0binary-hash
	var expectedBody bytes.Buffer
	// ディレクトリが先（"dir1/" < "file1.txt"）
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir1")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes()) // バイナリハッシュ
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "file1.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes())

	expectedData := []byte(fmt.Sprintf("tree %d\x00", expectedBody.Len()))
	expectedData = append(expectedData, expectedBody.Bytes()...)

	assert.Equal(t, expectedData, obj.Data)
}

func Test_TreeEntrySorting(t *testing.T) {
	tree := NewTree()
	tree.AddEntry(TreeEntry{Name: "b.txt", Hash: hash.SHA1{}, Mode: ModeFile})
	tree.AddEntry(TreeEntry{Name: "a.txt", Hash: hash.SHA1{}, Mode: ModeFile})
	tree.AddEntry(TreeEntry{Name: "dir1", Hash: hash.SHA1{}, Mode: ModeDir})
	tree.AddEntry(TreeEntry{Name: "dir2", Hash: hash.SHA1{}, Mode: ModeDir})

	obj := tree.serialize()

	// 期待されるソート順: a.txt, b.txt, dir1, dir2
	var expectedBody bytes.Buffer
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "a.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "b.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir1")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir2")
	expectedBody.WriteByte(0)
	expectedBody.Write(hash.SHA1{}.Bytes())

	expectedData := []byte(fmt.Sprintf("tree %d\x00", expectedBody.Len()))
	expectedData = append(expectedData, expectedBody.Bytes()...)

	assert.Equal(t, expectedData, obj.Data)
}
