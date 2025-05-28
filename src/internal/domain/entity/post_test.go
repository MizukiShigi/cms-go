package entity

import (
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestNewPost(t *testing.T) {
	userID := valueobject.NewUserID()
	title, _ := valueobject.NewPostTitle("テストタイトル")
	content, _ := valueobject.NewPostContent("テストコンテンツ")

	tests := []struct {
		name    string
		title   valueobject.PostTitle
		content valueobject.PostContent
		userID  valueobject.UserID
		status  valueobject.PostStatus
		wantErr bool
	}{
		{
			name:    "正常ケース: 下書き投稿作成",
			title:   title,
			content: content,
			userID:  userID,
			status:  valueobject.StatusDraft,
			wantErr: false,
		},
		{
			name:    "正常ケース: 公開投稿作成",
			title:   title,
			content: content,
			userID:  userID,
			status:  valueobject.StatusPublished,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(tt.title, tt.content, tt.userID, tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if post == nil {
				t.Error("投稿がnilです")
				return
			}

			// フィールドの検証
			if !post.Title.Equals(tt.title) {
				t.Errorf("Title = %v, want %v", post.Title, tt.title)
			}

			if !post.Content.Equals(tt.content) {
				t.Errorf("Content = %v, want %v", post.Content, tt.content)
			}

			if !post.UserID.Equals(tt.userID) {
				t.Errorf("UserID = %v, want %v", post.UserID, tt.userID)
			}

			if !post.Status.Equals(tt.status) {
				t.Errorf("Status = %v, want %v", post.Status, tt.status)
			}

			// IDが生成されていることを確認
			if post.ID.String() == "" {
				t.Error("投稿IDが生成されていません")
			}

			// タイムスタンプが設定されていることを確認
			if post.CreatedAt.IsZero() {
				t.Error("CreatedAtが設定されていません")
			}

			if post.UpdatedAt.IsZero() {
				t.Error("UpdatedAtが設定されていません")
			}

			// 公開投稿の場合FirstPublishedAtが設定されていることを確認
			if tt.status == valueobject.StatusPublished {
				if post.FirstPublishedAt == nil || post.FirstPublishedAt.IsZero() {
					t.Error("公開投稿のFirstPublishedAtが設定されていません")
				}
			}

			// ContentUpdatedAtが設定されていることを確認
			if post.ContentUpdatedAt == nil || post.ContentUpdatedAt.IsZero() {
				t.Error("ContentUpdatedAtが設定されていません")
			}
		})
	}
}

func TestPost_AddTag(t *testing.T) {
	userID := valueobject.NewUserID()
	title, _ := valueobject.NewPostTitle("テストタイトル")
	content, _ := valueobject.NewPostContent("テストコンテンツ")
	post, _ := NewPost(title, content, userID, valueobject.StatusDraft)

	tag1, _ := valueobject.NewTagName("golang")
	tag2, _ := valueobject.NewTagName("testing")

	tests := []struct {
		name        string
		tag         valueobject.TagName
		setup       func()
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: タグ追加",
			tag:     tag1,
			setup:   func() {},
			wantErr: false,
		},
		{
			name:    "正常ケース: 2番目のタグ追加",
			tag:     tag2,
			setup:   func() {},
			wantErr: false,
		},
		{
			name:        "異常ケース: 重複タグ追加",
			tag:         tag1,
			setup:       func() {},
			wantErr:     true,
			expectedErr: "Tag already exists",
		},
		{
			name: "異常ケース: 最大タグ数超過",
			tag:  tag1,
			setup: func() {
				// 投稿を新規作成して10個のタグを追加
				post, _ = NewPost(title, content, userID, valueobject.StatusDraft)
				for i := 0; i < 10; i++ {
					tagName, _ := valueobject.NewTagName("tag" + string(rune('0'+i)))
					post.AddTag(tagName)
				}
			},
			wantErr:     true,
			expectedErr: "Maximum number of tags reached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := post.AddTag(tt.tag)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("期待されたエラーメッセージ = %v, 実際のエラーメッセージ = %v", tt.expectedErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
		})
	}
}

func TestPost_SetStatus(t *testing.T) {
	userID := valueobject.NewUserID()
	title, _ := valueobject.NewPostTitle("テストタイトル")
	content, _ := valueobject.NewPostContent("テストコンテンツ")

	tests := []struct {
		name         string
		initialStatus valueobject.PostStatus
		targetStatus  valueobject.PostStatus
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "正常ケース: 下書きから公開",
			initialStatus: valueobject.StatusDraft,
			targetStatus:  valueobject.StatusPublished,
			wantErr:      false,
		},
		{
			name:         "正常ケース: 下書きから削除",
			initialStatus: valueobject.StatusDraft,
			targetStatus:  valueobject.StatusDeleted,
			wantErr:      false,
		},
		{
			name:         "正常ケース: 公開から非公開",
			initialStatus: valueobject.StatusPublished,
			targetStatus:  valueobject.StatusPrivate,
			wantErr:      false,
		},
		{
			name:         "正常ケース: 非公開から公開",
			initialStatus: valueobject.StatusPrivate,
			targetStatus:  valueobject.StatusPublished,
			wantErr:      false,
		},
		{
			name:         "正常ケース: 非公開から削除",
			initialStatus: valueobject.StatusPrivate,
			targetStatus:  valueobject.StatusDeleted,
			wantErr:      false,
		},
		{
			name:         "異常ケース: 公開から削除（不正な遷移）",
			initialStatus: valueobject.StatusPublished,
			targetStatus:  valueobject.StatusDeleted,
			wantErr:      true,
			expectedErr:  "Only draft and private posts can be deleted",
		},
		{
			name:         "異常ケース: 削除から公開（不正な遷移）",
			initialStatus: valueobject.StatusDeleted,
			targetStatus:  valueobject.StatusPublished,
			wantErr:      true,
			expectedErr:  "Only draft and private posts can be published",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, _ := NewPost(title, content, userID, tt.initialStatus)
			
			// 初期状態が公開以外で公開状態に設定したい場合は一度公開状態にする
			if tt.initialStatus == valueobject.StatusPrivate && post.Status != valueobject.StatusPrivate {
				post.SetStatus(valueobject.StatusPublished)
				post.SetStatus(valueobject.StatusPrivate)
			}
			if tt.initialStatus == valueobject.StatusDeleted && post.Status != valueobject.StatusDeleted {
				post.SetStatus(valueobject.StatusDeleted)
			}

			err := post.SetStatus(tt.targetStatus)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("期待されたエラーメッセージ = %v, 実際のエラーメッセージ = %v", tt.expectedErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if !post.Status.Equals(tt.targetStatus) {
				t.Errorf("Status = %v, want %v", post.Status, tt.targetStatus)
			}
		})
	}
}

func TestParsePost(t *testing.T) {
	id, _ := valueobject.ParsePostID("550e8400-e29b-41d4-a716-446655440000")
	title, _ := valueobject.NewPostTitle("テストタイトル")
	content, _ := valueobject.NewPostContent("テストコンテンツ")
	userID := valueobject.NewUserID()
	status := valueobject.StatusPublished
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()
	firstPublishedAt := time.Now().Add(-12 * time.Hour)
	contentUpdatedAt := time.Now().Add(-6 * time.Hour)
	tag1, _ := valueobject.NewTagName("golang")
	tag2, _ := valueobject.NewTagName("testing")
	tags := []valueobject.TagName{tag1, tag2}

	post := ParsePost(
		id,
		title,
		content,
		userID,
		status,
		createdAt,
		updatedAt,
		&firstPublishedAt,
		&contentUpdatedAt,
		tags,
	)

	if post == nil {
		t.Fatal("投稿がnilです")
	}

	// 各フィールドの検証
	if !post.ID.Equals(id) {
		t.Errorf("ID = %v, want %v", post.ID, id)
	}

	if !post.Title.Equals(title) {
		t.Errorf("Title = %v, want %v", post.Title, title)
	}

	if !post.Content.Equals(content) {
		t.Errorf("Content = %v, want %v", post.Content, content)
	}

	if !post.UserID.Equals(userID) {
		t.Errorf("UserID = %v, want %v", post.UserID, userID)
	}

	if !post.Status.Equals(status) {
		t.Errorf("Status = %v, want %v", post.Status, status)
	}

	if !post.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", post.CreatedAt, createdAt)
	}

	if !post.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", post.UpdatedAt, updatedAt)
	}

	if len(post.Tags) != 2 {
		t.Errorf("Tags length = %v, want %v", len(post.Tags), 2)
	}
}