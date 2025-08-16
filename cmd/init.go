package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func initRepository(targetDir string) error {
	// 絶対パスに変換
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// .pitディレクトリのパス
	pitDir := filepath.Join(absPath, ".pit")

	// .pitディレクトリが既に存在するかチェック
	if _, err := os.Stat(pitDir); !os.IsNotExist(err) {
		return fmt.Errorf("repository already exists at %s", pitDir)
	}

	// ディレクトリ構造を作成
	if err := createRepositoryStructure(pitDir); err != nil {
		return fmt.Errorf("failed to create repository structure: %w", err)
	}

	// 初期ファイルを作成
	if err := createInitialFiles(pitDir); err != nil {
		return fmt.Errorf("failed to create initial files: %w", err)
	}

	fmt.Printf("Initialized empty Pit repository in %s\n", pitDir)
	return nil
}

func createRepositoryStructure(pitDir string) error {
	// 必要なディレクトリを作成
	dirs := []string{
		pitDir,
		filepath.Join(pitDir, "objects"),
		filepath.Join(pitDir, "refs"),
		filepath.Join(pitDir, "refs", "heads"),
		filepath.Join(pitDir, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func createInitialFiles(pitDir string) error {
	// HEAD ファイルを作成 (デフォルトブランチはmain)
	headPath := filepath.Join(pitDir, "HEAD")
	headContent := "ref: refs/heads/main\n"
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %w", err)
	}

	// config ファイルを作成
	configPath := filepath.Join(pitDir, "config")
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// description ファイルを作成
	descPath := filepath.Join(pitDir, "description")
	descContent := "Unnamed repository; edit this file 'description' to name the repository.\n"
	if err := os.WriteFile(descPath, []byte(descContent), 0644); err != nil {
		return fmt.Errorf("failed to create description file: %w", err)
	}

	return nil
}

// Kong version of init command
type InitCmd struct {
	Directory string `arg:"" optional:"" help:"Target directory for initialization (default: current directory)"`
}

func (cmd *InitCmd) Run() error {
	targetDir := "."
	if cmd.Directory != "" {
		targetDir = cmd.Directory
	}
	return initRepository(targetDir)
}
