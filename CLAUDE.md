# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# CMSプロジェクト ドキュメント

## 原則
- 日本語で回答してください

## プロジェクト概要
このプロジェクトはGo言語を使用したCMS（Content Management System）のバックエンドシステムです。

## プロジェクト構造
```
src/
├── cmd/            # アプリケーションのエントリーポイント
├── config/         # 設定ファイル
├── infrastructure/ # インフラストラクチャ層（DB、外部サービス連携など）
├── internal/       # 内部パッケージ（ビジネスロジック、ドメインモデルなど）
│   ├── domain/     # ドメイン層
│   │   ├── entity/      # ドメインエンティティ
│   │   ├── valueobject/ # 値オブジェクト
│   │   ├── repository/  # リポジトリインターフェース
│   │   ├── service/     # ドメインサービス
│   │   └── context/     # ドメインコンテキスト
│   ├── usecase/    # ユースケース層
│   └── presentation/ # プレゼンテーション層
│       ├── controller/ # HTTPコントローラー
│       ├── middleware/ # HTTPミドルウェア
│       └── helper/     # プレゼンテーション層のヘルパー
├── migrations/     # データベースマイグレーションファイル
└── tmp/           # 一時ファイル
```

## アーキテクチャ

### オニオンアーキテクチャ
このプロジェクトはオニオンアーキテクチャを採用しています。オニオンアーキテクチャは、依存関係を内側に向かって流れるように設計されており、以下の特徴があります：

1. **依存関係の方向**
   - 内側の層は外側の層に依存しない
   - 依存関係は常に内側に向かって流れる
   - 外側の層は内側の層のインターフェースに依存

2. **層の構成（内側から外側）**
   ```
   [ドメインモデル]
        ↑
   [ドメインサービス]
        ↑
   [ユースケース]
        ↑
   [インターフェースアダプター]
        ↑
   [フレームワークとドライバー]
   ```

3. **各層の役割**
   - ドメインモデル: ビジネスロジックの中核
   - ドメインサービス: ドメインの操作とルール
   - ユースケース: アプリケーションのユースケース
   - インターフェースアダプター: 外部との接続
   - フレームワークとドライバー: 技術的な実装

4. **利点**
   - ビジネスロジックの独立性
   - テスタビリティの向上
   - フレームワークやデータベースの交換が容易
   - ドメインの純粋性の保持

## アーキテクチャ詳細

### ドメイン層 (internal/domain)
ドメイン層はビジネスロジックの中核を担う層です。

#### エンティティ (entity)
- ビジネスオブジェクトの定義
- ドメインの中心となるオブジェクト
- 例：User, Article, Category など

#### 値オブジェクト (valueobject)
- 不変な値の表現
- ドメイン固有の値の型定義
- 例：Email, Password, Title など

#### リポジトリインターフェース (repository)
- データアクセスの抽象化
- ドメインオブジェクトの永続化方法の定義
- 例：UserRepository, ArticleRepository など

#### ドメインサービス (service)
- 複数のエンティティにまたがるビジネスロジック
- エンティティに属さない操作の実装
- 例：認証サービス、検索サービス など

#### ドメインコンテキスト (context)
- ドメインの実行コンテキスト
- トランザクション管理
- ドメインイベントの管理

### ユースケース層 (internal/usecase)
- アプリケーションのユースケースの実装
- ドメイン層とプレゼンテーション層の仲介
- トランザクション管理
- 入力値のバリデーション
- ドメインオブジェクトの操作

### プレゼンテーション層 (internal/presentation)
HTTPリクエストの処理とレスポンスの生成を担当します。

#### コントローラー (controller)
- HTTPリクエストの受付
- リクエストパラメータのバリデーション
- ユースケースの呼び出し
- レスポンスの生成

#### ミドルウェア (middleware)
- 認証・認可
- リクエストロギング
- エラーハンドリング
- CORS設定
- レート制限

#### ヘルパー (helper)
- レスポンス形式の統一
- エラーメッセージの管理
- 共通ユーティリティ関数

### インフラストラクチャ層 (infrastructure)
外部システムとの統合とデータアクセスを担当します。

#### repository
- ドメインリポジトリインターフェースの実装
- SQLBoilerを使用したデータベースアクセス
- トランザクション管理の実装

#### service
- 外部サービスとの統合（JWT生成、Google Cloud Storage等）
- ドメインサービスインターフェースの実装

#### logger
- 構造化ログ（slog）の設定とハンドラー

## 重要な実装パターン

### 依存性注入
- インターフェースによる依存関係の抽象化
- テスト時のモック注入
- 各層の独立性確保

### エラーハンドリング
- カスタムエラー型による分類（`internal/domain/valueobject/error.go`）
- 適切なHTTPステータスコードへのマッピング
- 日本語エラーメッセージの提供

### 認証・認可フロー
1. JWTトークンの生成（ログイン時）
2. ミドルウェアでのトークン検証
3. ユーザーコンテキストの設定
4. 認証が必要なエンドポイントでの利用

## 主要な依存関係
- Web Framework: `github.com/gorilla/mux`
- データベース: PostgreSQL (`github.com/lib/pq`)
- ORM: SQLBoiler (`github.com/volatiletech/sqlboiler/v4`)
- 設定管理: Viper (`github.com/spf13/viper`)
- バリデーション: `github.com/go-playground/validator/v10`
- JWT認証: `github.com/golang-jwt/jwt/v4`

## コーディング規約
1. パッケージ構成
   - `internal/`: アプリケーションの内部ロジック
   - `infrastructure/`: 外部サービスとの連携
   - `cmd/`: アプリケーションのエントリーポイント

2. コメント規約
   - パッケージ、関数、構造体には日本語でコメントを記載
   - 複雑なロジックには詳細な説明を追加
   - TODOコメントは日本語で記載

3. エラーハンドリング
   - エラーメッセージは日本語で返却
   - エラーの種類に応じて適切なHTTPステータスコードを設定

## 開発環境
- Go 1.24.0
- Air (ホットリロード)
- Docker
- SQLBoiler (ORM)

## 開発コマンド

### テスト実行
```bash
# 全テスト実行（SQLBoilerモデルのテストを除く）
make test

# 単一テスト実行
cd src && go test ./internal/usecase/login_user_test.go
cd src && go test ./internal/domain/entity/user_test.go

# 特定パッケージのテスト
cd src && go test ./internal/usecase/
cd src && go test ./internal/domain/entity/
```

### コード品質
```bash
# コードフォーマット
make fmt

# コード静的解析
make vet
```

### データベース
```bash
# SQLBoilerモデル生成（マイグレーション）
make migration

# PostgreSQL接続設定: sqlboiler.toml
# デフォルト: host=localhost, port=5432, dbname=cms, user=postgres, pass=password
```

### アプリケーション実行
```bash
# APIサーバー起動
cd src && go run cmd/api/main.go

# Docker環境での実行
docker-compose up
```

### Google Cloud Platform
```bash
# GCP認証
make login

# Terraform操作
make plan    # 実行計画確認
make apply   # インフラ適用
make destroy # インフラ削除

# Docker操作
make docker-build  # イメージビルド
make docker-push   # イメージプッシュ
```

## デプロイメント
- Dockerコンテナ化対応
- 環境変数による設定管理
- マイグレーション自動実行

## セキュリティ考慮事項
- JWTによる認証
- パスワードのハッシュ化
- 環境変数による機密情報管理

## テスト戦略

### テストの優先順位
以下の順序でテストを実装することを推奨します：

1. **ドメイン層のテスト（最優先）**
   - エンティティの振る舞い
   - 値オブジェクトのバリデーション
   - ドメインサービスのロジック
   - 理由：ビジネスロジックの中核であり、最も重要な部分

2. **ユースケース層のテスト（次に優先）**
   - ユースケースの正常系
   - エラーケースの処理
   - トランザクション管理
   - 理由：アプリケーションの主要な機能をカバー

3. **プレゼンテーション層のテスト（必要に応じて）**
   - コントローラーの入力バリデーション
   - レスポンス形式の検証
   - 理由：HTTPリクエスト/レスポンスの処理は比較的シンプル

### テストの種類と実装方針

#### 1. ユニットテスト
- **対象**: ドメイン層、ユースケース層
- **実装方針**:
  - モックを使用して依存関係を分離
  - ビジネスロジックの検証に焦点
  - エッジケースのカバー

#### 2. インテグレーションテスト
- **対象**: リポジトリの実装、ユースケース
- **実装方針**:
  - 実際のデータベースを使用
  - トランザクションの検証
  - データの整合性確認

#### 3. E2Eテスト
- **対象**: 主要なユースケース
- **実装方針**:
  - 実際のHTTPリクエストを使用
  - 認証フローの検証
  - 重要なビジネスフローの確認

### テストコードの配置
```
src/
├── internal/
│   ├── domain/
│   │   └── entity/
│   │       ├── user.go
│   │       └── user_test.go    # ドメインテスト
│   ├── usecase/
│   │   └── user/
│   │       ├── create.go
│   │       └── create_test.go  # ユースケーステスト
│   └── presentation/
│       └── controller/
│           ├── user.go
│           └── user_test.go    # コントローラーテスト
└── tests/
    ├── integration/            # インテグレーションテスト
    └── e2e/                    # E2Eテスト
```

### テストカバレッジの目標
- ドメイン層: 90%以上
- ユースケース層: 80%以上
- プレゼンテーション層: 70%以上

### テストデータ管理
- テストデータは各テストファイル内で直接定義
- フィクスチャは未導入（今後の課題）
- テストデータは環境に依存しない形で管理

## パフォーマンス最適化
- データベースインデックスの適切な設定
- キャッシュの活用
- N+1問題の回避

## 監視・ロギング
- エラーログの適切な記録
- パフォーマンスメトリクスの収集
- アクセスログの記録

## 利用可能なコマンド
- `/project:list-issue` - issueをリストで取得
- `/project:fix-issue` - issueを解決

# important-instruction-reminders
Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
