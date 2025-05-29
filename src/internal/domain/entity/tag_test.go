package entity

import (
	"testing"
	"time"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestNewTagWithName(t *testing.T) {
	tagName, err := valueobject.NewTagName("golang")
	if err != nil {
		t.Fatalf("タグ名作成に失敗: %v", err)
	}

	tag := NewTagWithName(tagName)

	if tag == nil {
		t.Fatal("タグがnilです")
	}

	// フィールドの検証
	if !tag.Name.Equals(tagName) {
		t.Errorf("Name = %v, want %v", tag.Name, tagName)
	}

	// タイムスタンプが設定されていることを確認
	if tag.CreatedAt.IsZero() {
		t.Error("CreatedAtが設定されていません")
	}

	if tag.UpdatedAt.IsZero() {
		t.Error("UpdatedAtが設定されていません")
	}

	// CreatedAtとUpdatedAtが同じ時刻であることを確認（新規作成のため）
	if !tag.CreatedAt.Equal(tag.UpdatedAt) {
		t.Error("新規作成時はCreatedAtとUpdatedAtが同じ時刻である必要があります")
	}
}

func TestParseTag(t *testing.T) {
	tagID := valueobject.NewTagID()

	tagName, err := valueobject.NewTagName("testing")
	if err != nil {
		t.Fatalf("タグ名作成に失敗: %v", err)
	}

	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	tag := ParseTag(tagID, tagName, createdAt, updatedAt)

	if tag == nil {
		t.Fatal("タグがnilです")
	}

	// 各フィールドの検証
	if !tag.ID.Equals(tagID) {
		t.Errorf("ID = %v, want %v", tag.ID, tagID)
	}

	if !tag.Name.Equals(tagName) {
		t.Errorf("Name = %v, want %v", tag.Name, tagName)
	}

	if !tag.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", tag.CreatedAt, createdAt)
	}

	if !tag.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", tag.UpdatedAt, updatedAt)
	}
}

func TestTag_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 有効なタグ名",
			tagName: "golang",
			wantErr: false,
		},
		{
			name:    "正常ケース: 数字を含むタグ名",
			tagName: "go1-19",
			wantErr: false,
		},
		{
			name:    "正常ケース: アンダースコアを含むタグ名",
			tagName: "web_development",
			wantErr: false,
		},
		{
			name:        "異常ケース: 長すぎるタグ名",
			tagName:     "this-is-a-very-long-tag-name-that-exceeds-the-maximum-allowed-length",
			wantErr:     true,
			expectedErr: "TagName is too long",
		},
		{
			name:        "異常ケース: 特殊文字を含むタグ名",
			tagName:     "go@lang",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagName, err := valueobject.NewTagName(tt.tagName)

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

			// タグエンティティの作成
			tag := NewTagWithName(tagName)
			if tag == nil {
				t.Error("タグがnilです")
			}
		})
	}
}
