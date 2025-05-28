package valueobject

import (
	"testing"
)

func TestNewPostStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		wantErr     bool
		expectedErr string
		expected    PostStatus
	}{
		{
			name:     "正常ケース: draft",
			status:   "draft",
			wantErr:  false,
			expected: StatusDraft,
		},
		{
			name:     "正常ケース: published",
			status:   "published",
			wantErr:  false,
			expected: StatusPublished,
		},
		{
			name:     "正常ケース: private",
			status:   "private",
			wantErr:  false,
			expected: StatusPrivate,
		},
		{
			name:     "正常ケース: deleted",
			status:   "deleted",
			wantErr:  false,
			expected: StatusDeleted,
		},
		{
			name:        "異常ケース: 無効なステータス",
			status:      "invalid",
			wantErr:     true,
			expectedErr: "Invalid post status",
		},
		{
			name:        "異常ケース: 空文字",
			status:      "",
			wantErr:     true,
			expectedErr: "Invalid post status",
		},
		{
			name:        "異常ケース: 大文字",
			status:      "DRAFT",
			wantErr:     true,
			expectedErr: "Invalid post status",
		},
		{
			name:        "異常ケース: 混合文字",
			status:      "Draft",
			wantErr:     true,
			expectedErr: "Invalid post status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postStatus, err := NewPostStatus(tt.status)

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

			if !postStatus.Equals(tt.expected) {
				t.Errorf("NewPostStatus() = %v, want %v", postStatus, tt.expected)
			}
		})
	}
}

func TestPostStatus_Constants(t *testing.T) {
	// 定数の値を確認
	tests := []struct {
		name     string
		constant PostStatus
		expected string
	}{
		{
			name:     "StatusDraft",
			constant: StatusDraft,
			expected: "draft",
		},
		{
			name:     "StatusPublished",
			constant: StatusPublished,
			expected: "published",
		},
		{
			name:     "StatusPrivate",
			constant: StatusPrivate,
			expected: "private",
		},
		{
			name:     "StatusDeleted",
			constant: StatusDeleted,
			expected: "deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant.String() != tt.expected {
				t.Errorf("%s.String() = %v, want %v", tt.name, tt.constant.String(), tt.expected)
			}
		})
	}
}

func TestPostStatus_String(t *testing.T) {
	status := StatusPublished
	expected := "published"

	if status.String() != expected {
		t.Errorf("String() = %v, want %v", status.String(), expected)
	}
}

func TestPostStatus_Equals(t *testing.T) {
	status1 := StatusDraft
	status2 := StatusDraft
	status3 := StatusPublished

	// 同じステータス
	if !status1.Equals(status2) {
		t.Error("同じステータス同士の比較がfalseになりました")
	}

	// 異なるステータス
	if status1.Equals(status3) {
		t.Error("異なるステータス同士の比較がtrueになりました")
	}

	// 自身との比較
	if !status1.Equals(status1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestPostStatus_AllValues(t *testing.T) {
	// 全てのステータス値に対してNewPostStatusが正常に動作することを確認
	statuses := []string{"draft", "published", "private", "deleted"}

	for _, status := range statuses {
		t.Run("status_"+status, func(t *testing.T) {
			postStatus, err := NewPostStatus(status)
			if err != nil {
				t.Errorf("有効なステータス '%s' でエラーが発生しました: %v", status, err)
			}

			if postStatus.String() != status {
				t.Errorf("String() = %v, want %v", postStatus.String(), status)
			}
		})
	}
}