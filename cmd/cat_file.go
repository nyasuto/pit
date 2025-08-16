package cmd

import (
	"fmt"

	"github.com/nyasuto/pit/internal/objects"
)

// cat-file command
type CatFileCmd struct {
	Print bool   `short:"p" help:"Print the object from the .pit/objects directory"`
	Type  bool   `short:"t" help:"Print the type from the .pit/objects directory"`
	Hash  string `arg:"hash" help:"Hash of the file to print"`
}

func (cmd *CatFileCmd) Validate() error {
	if cmd.Hash == "" {
		return fmt.Errorf("hash must be specified")
	}
	if cmd.Print && cmd.Type {
		return fmt.Errorf("cannot specify both -p and -t")
	}
	return nil
}

func (cmd *CatFileCmd) Run() error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	var err error

	obj, err := objects.ReadFromHash(cmd.Hash)

	if err != nil {
		return fmt.Errorf("failed to read object %q: %w", cmd.Hash, err)
	}

	// デフォルト動作: オプション未指定時は -p 動作
	if !cmd.Print && !cmd.Type {
		cmd.Print = true
	}

	if cmd.Print {
		fmt.Print(obj.String())
		return nil
	}
	if cmd.Type {
		fmt.Println(obj.Type)
		return nil
	}
	return fmt.Errorf("no action specified, use -p to print or -t to print type")
}
