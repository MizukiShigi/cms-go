# GitHub Issue一覧取得

GitHubのissue一覧を取得して表示します。

## 機能
- open/closed/allの状態でフィルタリング
- ラベルでフィルタリング
- 日本語対応の見やすい表示

## 使用法
状態とラベルを指定してissue一覧を取得してください。

## 実行コマンド

```bash
# 全てのissueを表示
gh issue list --json number,title,state,labels,assignees,createdAt,updatedAt | jq -r '.[] | "ID: \(.number | tostring) | \(.state | ascii_upcase) | \(.title) | 作成: \(.createdAt[:10]) | 更新: \(.updatedAt[:10]) | ラベル: \(if .labels | length > 0 then (.labels | map(.name) | join(", ")) else "なし" end)"'
```

パラメータを指定して、適切なコマンドを実行してください。