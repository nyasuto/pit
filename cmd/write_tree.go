package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nyasuto/pit/internal/objects"
)

// cat-file command
type WriteTreeCmd struct {
}

func (cmd *WriteTreeCmd) Validate() error {

	return nil
}

func (cmd *WriteTreeCmd) Run() error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	tree, err := buildTreeFromDirectory(".")
	if err != nil {
		return err
	}

	treeObj := tree.Serialize()
	_, err = objects.Write(treeObj)
	if err != nil {
		return err
	}

	fmt.Println(treeObj.Hash.String())

	return nil

}

func shouldIgnore(name string) bool {
	return name == ".pit" || name == ".git" // シンプルに.pitを無視
}

func buildTreeFromDirectory(path string) (*objects.Tree, error) {
	tree := objects.NewTree()

	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		if shouldIgnore(entry.Name()) {
			continue // .pitをスキップ
		}
		if entry.IsDir() {
			// 1. サブディレクトリを再帰処理
			subTree, _ := buildTreeFromDirectory(filepath.Join(path, entry.Name()))

			// 2. サブtreeオブジェクトを作成・保存
			subTreeObj := subTree.Serialize()
			objects.Write(subTreeObj) // ← ここで保存

			// 3. 親treeにサブtreeのハッシュを追加
			tree.AddEntry(objects.TreeEntry{
				Mode: objects.ModeDir,
				Name: entry.Name(),
				Hash: subTreeObj.Hash,
			})

		} else {
			// ファイルもblobとして保存
			data, _ := os.ReadFile(filepath.Join(path, entry.Name()))
			blob := objects.NewBlob(data)
			objects.Write(blob) // ← ここで保存

			tree.AddEntry(objects.TreeEntry{
				Mode: objects.ModeFile,
				Name: entry.Name(),
				Hash: blob.Hash,
			})
		}
	}

	return tree, nil
}
