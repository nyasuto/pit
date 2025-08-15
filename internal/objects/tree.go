package objects

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/nyasuto/pit/pkg/hash"
)

type TreeEntry struct {
	Name string     // Name of the entry (file or sub-tree)
	Hash hash.SHA1  // SHA1 hash of the entry
	Mode ObjectMode // File mode (e.g., 0644 for files, 040000 for directories)
}
type Tree struct {
	Entries []TreeEntry // Entries in the tree
}

func isDirectory(mode ObjectMode) bool {
	return mode == ModeDir // 040000 in octal
}

func (t *Tree) toObject() object {
	return t.serialize()

}

func NewTree() *Tree {
	return &Tree{
		Entries: []TreeEntry{},
	}
}

func (t *Tree) AddEntry(entry TreeEntry) {
	t.Entries = append(t.Entries, entry)
}

// NewTree creates a new tree object from the provided entries.
func (t *Tree) serialize() object {
	entries := t.Entries
	// エントリのソート（Git仕様に準拠）
	sort.Slice(entries, func(i, j int) bool {
		nameI := entries[i].Name
		nameJ := entries[j].Name
		if isDirectory(entries[i].Mode) {
			nameI += "/"
		}
		if isDirectory(entries[j].Mode) {
			nameJ += "/"
		}
		return nameI < nameJ
	})

	var buf bytes.Buffer
	for _, entry := range entries {

		fmt.Fprintf(&buf, "%o %s", entry.Mode, entry.Name)
		buf.WriteByte(0) // Null byte to separate entries
		buf.Write(entry.Hash.Bytes())
	}

	return New(ObjectTypeTree, buf.Bytes())

}
