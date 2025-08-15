package objects

import (
	"fmt"
	"time"

	"github.com/nyasuto/pit/pkg/hash"
)

type Commit struct {
	Tree    hash.SHA1  // ルートTreeオブジェクトのハッシュ
	Parents *hash.SHA1 // 親コミットのハッシュ（マージ時は複数）
	Author  Person     // 作成者情報
	Message string     // コミットメッセージ
}

type Person struct {
	Name     string    // 名前
	Email    string    // メールアドレス
	When     time.Time // 日時
	TimeZone string    // タイムゾーン（例: "+0900"）
}

func (c *Commit) SetAuthor(name, email string) {
	c.Author = newPerson(name, email)
}

func newPerson(name, email string) Person {
	now := time.Now()
	return Person{
		Name:     name,
		Email:    email,
		When:     now,
		TimeZone: now.Format("-0700"), // Git形式のタイムゾーン (+0900, -0800等)
	}
}

func NewCommit(tree hash.SHA1, message string) *Commit {
	return NewCommitWithParent(tree, nil, message)
}

func NewCommitWithAuthor(tree hash.SHA1, message, authorName, authorEmail string) *Commit {
	commit := NewCommit(tree, message)
	commit.SetAuthor(authorName, authorEmail)
	return commit
}

func NewCommitWithParent(tree hash.SHA1, parents *hash.SHA1, message string) *Commit {
	return &Commit{
		Tree:    tree,
		Parents: parents,
		Message: message,
	}
}
func (c *Commit) Serialize() []byte {
	var data []byte
	
	// tree行
	data = append(data, []byte("tree "+c.Tree.String()+"\n")...)
	
	// parent行（存在する場合のみ）
	if c.Parents != nil {
		data = append(data, []byte("parent "+c.Parents.String()+"\n")...)
	}
	
	// author行とcommitter行（Author情報が設定されている場合）
	if c.Author.Name != "" {
		timestamp := c.Author.When.Unix()
		authorLine := fmt.Sprintf("author %s <%s> %d %s\n", 
			c.Author.Name, c.Author.Email, timestamp, c.Author.TimeZone)
		committerLine := fmt.Sprintf("committer %s <%s> %d %s\n", 
			c.Author.Name, c.Author.Email, timestamp, c.Author.TimeZone)
		
		data = append(data, []byte(authorLine)...)
		data = append(data, []byte(committerLine)...)
	}
	
	// 空行 + メッセージ
	data = append(data, []byte("\n"+c.Message+"\n")...)
	return data
}

func (c *Commit) ToObject() object {
	data := c.Serialize()
	return New(ObjectTypeCommit, data)
}
