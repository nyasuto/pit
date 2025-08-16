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

	testHash1 := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	testHash2 := hash.SHA1{0x0f, 0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00}
	tree := NewTree()
	tree.AddEntry(TreeEntry{Name: "file1.txt", Hash: testHash1, Mode: ModeFile})
	tree.AddEntry(TreeEntry{Name: "dir1", Hash: testHash2, Mode: ModeDir})

	obj := tree.Serialize()
	assert.Equal(t, ObjectTypeTree, obj.Type)

	// Gitの正しい形式: mode name\0binary-hash
	var expectedBody bytes.Buffer
	// ディレクトリが先（"dir1/" < "file1.txt"）
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir1")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash2.Bytes()) // バイナリハッシュ
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "file1.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash1.Bytes())

	expectedData := []byte(fmt.Sprintf("tree %d\x00", expectedBody.Len()))
	expectedData = append(expectedData, expectedBody.Bytes()...)

	assert.Equal(t, expectedData, obj.Data)
}

func Test_TreeEntrySorting(t *testing.T) {
	testHash1 := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	testHash2 := hash.SHA1{0x0f, 0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00}
	testHash3 := hash.SHA1{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	testHash4 := hash.SHA1{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f}

	tree := NewTree()
	tree.AddEntry(TreeEntry{Name: "b.txt", Hash: testHash1, Mode: ModeFile})
	tree.AddEntry(TreeEntry{Name: "a.txt", Hash: testHash2, Mode: ModeFile})
	tree.AddEntry(TreeEntry{Name: "dir1", Hash: testHash3, Mode: ModeDir})
	tree.AddEntry(TreeEntry{Name: "dir2", Hash: testHash4, Mode: ModeDir})

	obj := tree.Serialize()

	// 期待されるソート順: a.txt, b.txt, dir1, dir2
	var expectedBody bytes.Buffer
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "a.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash2.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeFile, "b.txt")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash1.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir1")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash3.Bytes())
	fmt.Fprintf(&expectedBody, "%o %s", ModeDir, "dir2")
	expectedBody.WriteByte(0)
	expectedBody.Write(testHash4.Bytes())

	expectedData := []byte(fmt.Sprintf("tree %d\x00", expectedBody.Len()))
	expectedData = append(expectedData, expectedBody.Bytes()...)

	assert.Equal(t, expectedData, obj.Data)
}

func Test_FindEntry(t *testing.T) {
	testHash1 := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	tree := NewTree()
	entry := TreeEntry{Name: "file.txt", Hash: testHash1, Mode: ModeFile}
	tree.AddEntry(entry)

	foundEntry, found := tree.FindEntry("file.txt")
	assert.True(t, found)
	assert.Equal(t, entry, foundEntry)

	_, notFound := tree.FindEntry("nonexistent.txt")
	assert.False(t, notFound)
}

func Test_RemoveEntry(t *testing.T) {
	testHash1 := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	tree := NewTree()
	entry := TreeEntry{Name: "file.txt", Hash: testHash1, Mode: ModeFile}
	tree.AddEntry(entry)

	removed := tree.RemoveEntry("file.txt")
	assert.True(t, removed)
	_, found := tree.FindEntry("file.txt")
	assert.False(t, found)

	notRemoved := tree.RemoveEntry("nonexistent.txt")
	assert.False(t, notRemoved)
}

func Test_UpdateEntry(t *testing.T) {
	testHash1 := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	tree := NewTree()
	entry := TreeEntry{Name: "file.txt", Hash: testHash1, Mode: ModeFile}
	tree.AddEntry(entry)

	newHash := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	err := tree.UpdateEntry("file.txt", newHash)
	assert.NoError(t, err)

	updatedEntry, found := tree.FindEntry("file.txt")
	assert.True(t, found)
	assert.Equal(t, newHash, updatedEntry.Hash)

	err = tree.UpdateEntry("nonexistent.txt", newHash)
	assert.Error(t, err)
}
