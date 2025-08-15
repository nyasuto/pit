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
func (t *Tree) FindEntry(name string) (TreeEntry, bool)       {
	for _, entry := range t.Entries {
		if entry.Name == name {
			return entry, true
		}
	}
	return TreeEntry{}, false
}
func (t *Tree) RemoveEntry(name string) bool                  {
	for i, entry := range t.Entries {
		if entry.Name == name {
			t.Entries = append(t.Entries[:i], t.Entries[i+1:]...)
			return true
		}
	}
	return false
}
func (t *Tree) UpdateEntry(name string, hash hash.SHA1) error {
	for i, entry := range t.Entries {
		if entry.Name == name {
			t.Entries[i].Hash = hash
			return nil
		}
	}
	return fmt.Errorf("entry %s not found", name)
}
func NewTree() *Tree {
	return &Tree{
		Entries: []TreeEntry{},
	}
}

func (t *Tree) AddEntry(entry TreeEntry) error{
	if entry.Name == "" {
		return fmt.Errorf("entry name cannot be empty")
	}
	if entry.Hash == (hash.SHA1{}) {
		return fmt.Errorf("entry hash cannot be empty")
	}
	if entry.Mode == 0 {
		return fmt.Errorf("entry mode cannot be zero")
	}
	// Check for duplicate entries
	_, find := t.FindEntry(entry.Name)
	if find {
		return fmt.Errorf("entry %s already exists", entry.Name)
	}
	t.Entries = append(t.Entries, entry)
	return nil
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
