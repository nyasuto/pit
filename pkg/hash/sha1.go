package hash

import (
	"crypto/sha1"
	"fmt"
)

type SHA1 string

func Hash(data []byte) SHA1 {
	// SHA1の計算を行うために、crypto/sha1パッケージを使用
	// ここでは、SHA1のハッシュ値を計算して返す関数を定義します

	// データが空の場合はゼロ値を返す
	if len(data) == 0 {
		return ""
	}
	hash := sha1.Sum(data)
	return SHA1(fmt.Sprintf("%x", hash))
}
