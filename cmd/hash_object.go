package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nyasuto/pit/internal/objects"
)

// hash-object command
type HashObjectCmd struct {
	Write bool   `short:"w" help:"Write the object to the .pit/objects directory"`
	Stdin bool   `help:"Read from stdin instead of a file"`
	File  string `arg:"" optional:"" help:"File to hash"`
}

func (cmd *HashObjectCmd) Validate() error {
	if cmd.Stdin && cmd.File != "" {
		return fmt.Errorf("cannot specify both --stdin and file argument")
	}
	if !cmd.Stdin && cmd.File == "" {
		return fmt.Errorf("requires exactly one file argument when not using --stdin")
	}
	return nil
}

func (cmd *HashObjectCmd) Run() error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	var data []byte
	var err error

	if cmd.Stdin {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		data, err = os.ReadFile(cmd.File)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", cmd.File, err)
		}
	}

	return cmd.processData(data)
}

func (cmd *HashObjectCmd) processData(data []byte) error {
	// blobオブジェクト作成
	blob := objects.NewBlob(data)

	// ハッシュ値取得
	hash := blob.Hash.String()

	// -w オプションが指定されていれば保存
	if cmd.Write {
		_, err := objects.Write(blob)
		if err != nil {
			return fmt.Errorf("failed to write object: %w", err)
		}
	}

	// ハッシュ値を出力
	fmt.Println(hash)
	return nil
}
