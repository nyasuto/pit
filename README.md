# 🕳️ Pit

**A tiny, educational Git implementation in Go**

Pitは、Gitの内部構造を学ぶために作る、小さくて可愛いGit実装です。

## 📝 プロジェクトの目的

- Gitの内部メカニズムを深く理解する
- Go言語でシンプルかつ読みやすい実装を作る
- 教育目的で、誰でも理解できるコードを書く
- 小さく始めて、徐々に機能を追加していく

## 🎯 コンセプト

```
"Small, Fast, and Educational"
```

本家Gitの複雑さを排除し、コア機能に集中

## 🏗️ アーキテクチャ

```
pit/
├── cmd/
│   └── pit/
│       └── main.go          # CLIエントリーポイント
├── internal/
│   ├── objects/             # Git Objects (Blob, Tree, Commit)
│   │   ├── blob.go
│   │   ├── tree.go
│   │   └── commit.go
│   ├── refs/                # References (branches, tags)
│   │   └── refs.go
│   ├── index/               # Staging area
│   │   └── index.go
│   └── storage/             # Object storage (.git/objects)
│       └── storage.go
├── pkg/
│   ├── hash/                # SHA-1 hashing
│   │   └── sha.go
│   └── compress/            # zlib compression
│       └── zlib.go
├── examples/                # 使用例
├── docs/                    # ドキュメント
│   └── internals.md        # 内部設計の解説
└── tests/                   # テストリポジトリ
```

## 🚀 開発ロードマップ

### Phase 1: Core Objects（基本実装）✨
**目標**: Gitの3つの基本オブジェクトを実装

- [x] プロジェクト初期設定
- [x] Blob Object（ファイル内容の保存）
- [x] Tree Object（ディレクトリ構造）
- [x] Commit Object（コミット情報）
- [x] SHA-1ハッシュ計算
- [x] zlib圧縮・展開

**コマンド**:
```bash
pit hash-object <file>    # ファイルをBlob Objectとして保存
pit cat-file <hash>       # オブジェクトの内容を表示
```

### Phase 2: Basic Commands（基本コマンド）🛠️
**目標**: 最小限のGit操作を可能に

- [ ] `pit init` - リポジトリ初期化
- [ ] `pit add` - ステージングエリアに追加
- [ ] `pit commit` - コミット作成
- [ ] `pit log` - コミット履歴表示
- [ ] `pit status` - 現在の状態表示

**この時点で**: 基本的なバージョン管理が可能に！

### Phase 3: Branches（ブランチ機能）🌿
**目標**: ブランチの作成と切り替え

- [ ] `pit branch` - ブランチ一覧・作成
- [ ] `pit checkout` - ブランチ切り替え
- [ ] HEAD参照の管理
- [ ] refs/heads/の実装

### Phase 4: Diff & Merge（差分とマージ）🔀
**目標**: 変更の可視化と統合

- [ ] `pit diff` - 差分表示（簡易版）
- [ ] `pit merge` - Fast-forwardマージのみ
- [ ] 3-way mergeの基礎実装（チャレンジ）

### Phase 5: Remote（リモート機能）🌍
**目標**: 他のPitリポジトリとの同期

- [ ] `pit clone` - ローカルクローン
- [ ] `pit push` - ローカルプッシュ
- [ ] `pit pull` - ローカルプル
- [ ] Packfile形式の実装（Optional）

### Phase 6: Performance（最適化）⚡
**目標**: 実用的な速度を実現

- [ ] オブジェクトキャッシュ
- [ ] インデックスの最適化
- [ ] 並行処理の導入
- [ ] ベンチマークテスト

## 🎓 学習ポイント

各フェーズで学べること：

1. **ファイルシステムとデータ構造**
   - Content-Addressable Storage
   - Merkle Tree構造

2. **圧縮アルゴリズム**
   - zlib圧縮の仕組み
   - データの効率的な保存

3. **ハッシュ関数**
   - SHA-1の役割
   - 整合性の保証

4. **グラフ理論**
   - DAG（有向非巡回グラフ）
   - コミットグラフの走査

## 🚦 Getting Started

```bash
# リポジトリのクローン
git clone https://github.com/nyasuto/pit.git

# ビルド
cd pit
go build -o pit cmd/pit/main.go

# 最初のコマンドを試す
./pit init my-repo
cd my-repo
echo "Hello, Pit!" > hello.txt
../pit add hello.txt
../pit commit -m "First commit with Pit!"
```

## 🧪 テスト

```bash
# ユニットテスト
go test ./...

# 統合テスト
./scripts/integration_test.sh

# ベンチマーク
go test -bench=. ./...
```

## 📊 進捗状況

| Phase | 状態 | 進捗 |
|-------|------|------|
| Phase 1: Core Objects | ✅ 完了 | 100% |
| Phase 2: Basic Commands | ⏳ 待機中 | 0% |
| Phase 3: Branches | ⏳ 待機中 | 0% |
| Phase 4: Diff & Merge | ⏳ 待機中 | 0% |
| Phase 5: Remote | ⏳ 待機中 | 0% |
| Phase 6: Performance | ⏳ 待機中 | 0% |

## 🤝 コントリビューション

個人学習プロジェクトですが、以下は大歓迎です：

- バグ報告
- コードレビュー
- アイデアの提案
- ドキュメントの改善

## 📚 参考資料

- [Pro Git - Git Internals](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
- [Write yourself a Git!](https://wyag.thb.lt/)
- [Building Git](https://shop.jcoglan.com/building-git/)

## 🏷️ バージョニング

- v0.1.0 - Phase 1 完了（Core Objects）
- v0.2.0 - Phase 2 完了（Basic Commands）
- v0.3.0 - Phase 3 完了（Branches）
- v1.0.0 - Phase 4 完了（基本機能完成）

## 📄 ライセンス

MIT License - 学習目的で自由に使用・改変してください

## 🕳️ なぜ「Pit」？

Git（ギット）の内部を掘り下げる「穴（Pit）」のように、深く学ぶプロジェクト：
- **深い**: Gitの内部構造を深く理解
- **掘る**: 表面的でなく、本質を掘り下げる
- **学ぶ**: 知識の穴を埋める教育ツール

---

**現在のステータス**: ✅ Phase 1 完了！Phase 2 開始準備中...

```go
// 完成したコア機能
✅ Blob:   NewBlob(data []byte) object
✅ Tree:   NewTree(), AddEntry(), serialize()  
✅ Commit: NewCommit(tree, message), SetAuthor()
✅ 全テスト成功、Git仕様準拠、互換性確保
```

*Last Updated: 2025-08-15*