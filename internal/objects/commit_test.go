package objects

import (
	"strings"
	"testing"

	"github.com/nyasuto/pit/pkg/hash"
)

func Test_NewCommit(t *testing.T) {
	treeHash := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	commit := NewCommit(treeHash, "Initial commit")

	if commit.Tree != treeHash {
		t.Errorf("Expected tree hash %s, got %s", treeHash.String(), commit.Tree.String())
	}

	if commit.Message != "Initial commit" {
		t.Errorf("Expected message 'Initial commit', got '%s'", commit.Message)
	}

	obj := commit.ToObject()
	if obj.Type != ObjectTypeCommit {
		t.Errorf("Expected object type %s, got %s", ObjectTypeCommit, obj.Type)
	}
}

func Test_NewCommitWithParent(t *testing.T) {
	treeHash := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	parentHash := hash.SHA1{0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}
	commit := NewCommitWithParent(treeHash, &parentHash, "Commit with parent")

	if commit.Tree != treeHash {
		t.Errorf("Expected tree hash %s, got %s", treeHash.String(), commit.Tree.String())
	}

	if commit.Parents == nil || *commit.Parents != parentHash {
		t.Errorf("Expected parent hash %s, got %v", parentHash.String(), commit.Parents)
	}

	if commit.Message != "Commit with parent" {
		t.Errorf("Expected message 'Commit with parent', got '%s'", commit.Message)
	}

	obj := commit.ToObject()
	if obj.Type != ObjectTypeCommit {
		t.Errorf("Expected object type %s, got %s", ObjectTypeCommit, obj.Type)
	}
}

func Test_SerializeCommit(t *testing.T) {
	treeHash := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	commit := NewCommit(treeHash, "Test commit")
	commit.SetAuthor("Test Author", "test@example.com")
	
	// シリアライズされたデータを取得
	data := commit.Serialize()
	dataStr := string(data)

	// 各行を個別に検証
	if !strings.Contains(dataStr, "tree "+treeHash.String()) {
		t.Error("Missing tree line")
	}
	
	if !strings.Contains(dataStr, "author Test Author <test@example.com>") {
		t.Error("Missing author line")
	}
	
	if !strings.Contains(dataStr, "committer Test Author <test@example.com>") {
		t.Error("Missing committer line")
	}
	
	if !strings.Contains(dataStr, "Test commit") {
		t.Error("Missing commit message")
	}
	
	// tree行が最初に来ることを確認
	if !strings.HasPrefix(dataStr, "tree "+treeHash.String()+"\n") {
		t.Error("Tree line should be first")
	}
}

func Test_NewCommitWithAuthor(t *testing.T) {
	treeHash := hash.SHA1{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	commit := NewCommitWithAuthor(treeHash, "Test message", "Jane Doe", "jane@example.com")
	
	if commit.Author.Name != "Jane Doe" {
		t.Errorf("Expected author name 'Jane Doe', got '%s'", commit.Author.Name)
	}
	
	if commit.Author.Email != "jane@example.com" {
		t.Errorf("Expected author email 'jane@example.com', got '%s'", commit.Author.Email)
	}
	
	if commit.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", commit.Message)
	}
}