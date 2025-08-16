package objects

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nyasuto/pit/pkg/hash"
)

type ObjectType string

const (
	ObjectTypeBlob   ObjectType = "blob"
	ObjectTypeTree   ObjectType = "tree"
	ObjectTypeCommit ObjectType = "commit"
)

type ObjectMode uint32

const (
	// Git標準ファイルモード
	ModeFile       ObjectMode = 0100644 // 通常ファイル
	ModeExecutable ObjectMode = 0100755 // 実行可能ファイル
	ModeSymlink    ObjectMode = 0120000 // シンボリックリンク
	ModeDir        ObjectMode = 0040000 // ディレクトリ
	ModeSubmodule  ObjectMode = 0160000 // サブモジュール
)

const gitObjectsDir = ".pit/objects"

type object struct {
	Type ObjectType // Type of the object (e.g., "blob", "tree", "commit")
	Hash hash.SHA1  // SHA1 hash of the object
	Data []byte     // Raw data of the object
}

func New(t ObjectType, data []byte) object {

	size := len(data)
	header := []byte(fmt.Sprintf("%s %d\x00", t, size))

	content := append(header, data...)
	h := hash.Hash(content)

	return object{
		Type: t,
		Hash: h,
		Data: content,
	}
}

func (o *object) String() string {
	return fmt.Sprintf("%s", o.Data)
}

func ReadFromHash(hashString string) (obj object, err error) {
	h, err := hash.Parse(hashString)
	if err != nil {
		return object{}, fmt.Errorf("invalid hash: %s", hashString)
	}
	path := filepath.Join(gitObjectsDir, h.String()[:2], h.String()[2:])
	return Read(path)
}

func Read(path string) (object, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return object{}, err
	}
	zr, err := zlib.NewReader(bytes.NewReader(raw))
	if err != nil {
		return object{}, err
	}
	defer zr.Close()

	inflated, err := io.ReadAll(zr)
	if err != nil {
		return object{}, err
	}
	// 先頭のヘッダーを解析
	headerEnd := bytes.IndexByte(inflated, 0)
	if headerEnd < 0 {
		return object{}, fmt.Errorf("invalid object data: no header found")
	}
	header := inflated[:headerEnd]
	data := inflated[headerEnd+1:]
	// ヘッダーからタイプとサイズを取得
	parts := bytes.SplitN(header, []byte(" "), 2)
	if len(parts) != 2 {
		return object{}, fmt.Errorf("invalid object header: %s", header)
	}
	t := ObjectType(parts[0])
	size := len(data)

	// サイズをヘッダーから取得
	sizeInHeader := 0
	fmt.Sscanf(string(parts[1]), "%d", &sizeInHeader)

	if size != sizeInHeader {
		return object{}, fmt.Errorf("object size mismatch: expected %d, got %d", size, len(data))
	}
	h := hash.Hash(inflated)
	return object{
		Type: t,
		Hash: h,
		Data: inflated,
	}, nil
}

func Write(o object) (name string, err error) {
	if o.Type != ObjectTypeBlob && o.Type != ObjectTypeTree && o.Type != ObjectTypeCommit {
		return "", fmt.Errorf("unsupported object type: %s", o.Type)
	}
	hex := o.Hash.String()
	if len(hex) < 3 {
		return "", fmt.Errorf("invalid hash: %q", hex)
	}
	dir := filepath.Join(gitObjectsDir, hex[:2])
	path := filepath.Join(dir, hex[2:])
	// ディレクトリ作成
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	// zlib で圧縮
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write([]byte(o.Data)); err != nil {
		_ = zw.Close()
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}

	// Git は0444で置くことが多い（読み取り専用）
	if err := os.WriteFile(path, buf.Bytes(), 0o444); err != nil {
		return "", err
	}

	return path, err
}
