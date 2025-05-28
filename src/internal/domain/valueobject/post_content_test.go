package valueobject

import (
	"strings"
	"testing"
)

func TestNewPostContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 空のコンテンツ",
			content: "",
			wantErr: false,
		},
		{
			name:    "正常ケース: 短いコンテンツ",
			content: "これはテストコンテンツです。",
			wantErr: false,
		},
		{
			name:    "正常ケース: 通常の長さのコンテンツ",
			content: strings.Repeat("テスト", 100),
			wantErr: false,
		},
		{
			name:    "正常ケース: 最大長のコンテンツ（10000文字）",
			content: strings.Repeat("a", 10000),
			wantErr: false,
		},
		{
			name:    "正常ケース: 改行を含むコンテンツ",
			content: "第1行目\n第2行目\n第3行目",
			wantErr: false,
		},
		{
			name:    "正常ケース: 特殊文字を含むコンテンツ",
			content: "特殊文字：!@#$%^&*()_+-=[]{}|;':\",./<>?",
			wantErr: false,
		},
		{
			name:    "正常ケース: Unicode文字を含むコンテンツ",
			content: "絵文字：😀😃😄😁😆😅🤣😂🙂🙃😉😊",
			wantErr: false,
		},
		{
			name:        "異常ケース: 長すぎるコンテンツ（10001文字）",
			content:     strings.Repeat("a", 10001),
			wantErr:     true,
			expectedErr: "Content is too long",
		},
		{
			name:        "異常ケース: 非常に長いコンテンツ",
			content:     strings.Repeat("テスト", 5000),
			wantErr:     true,
			expectedErr: "Content is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := NewPostContent(tt.content)

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

			if content.String() != tt.content {
				t.Errorf("PostContent.String() = %v, want %v", content.String(), tt.content)
			}
		})
	}
}

func TestPostContent_String(t *testing.T) {
	content := "テストコンテンツ"
	postContent, err := NewPostContent(content)
	if err != nil {
		t.Fatalf("PostContent作成に失敗: %v", err)
	}

	if postContent.String() != content {
		t.Errorf("String() = %v, want %v", postContent.String(), content)
	}
}

func TestPostContent_Equals(t *testing.T) {
	content1, _ := NewPostContent("同じコンテンツ")
	content2, _ := NewPostContent("同じコンテンツ")
	content3, _ := NewPostContent("異なるコンテンツ")

	// 同じコンテンツ
	if !content1.Equals(content2) {
		t.Error("同じコンテンツ同士の比較がfalseになりました")
	}

	// 異なるコンテンツ
	if content1.Equals(content3) {
		t.Error("異なるコンテンツ同士の比較がtrueになりました")
	}

	// 自身との比較
	if !content1.Equals(content1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestPostContent_Boundary(t *testing.T) {
	// 境界値テスト：10000文字ちょうど
	content10000 := strings.Repeat("a", 10000)
	postContent, err := NewPostContent(content10000)
	if err != nil {
		t.Errorf("10000文字のコンテンツでエラーが発生しました: %v", err)
	}
	if len(postContent.String()) != 10000 {
		t.Errorf("10000文字のコンテンツ長 = %d, want 10000", len(postContent.String()))
	}

	// 境界値テスト：10001文字
	content10001 := strings.Repeat("a", 10001)
	_, err = NewPostContent(content10001)
	if err == nil {
		t.Error("10001文字のコンテンツでエラーが発生しませんでした")
	}
}

func TestPostContent_WithEmptyString(t *testing.T) {
	// 空文字列の場合
	postContent, err := NewPostContent("")
	if err != nil {
		t.Errorf("空文字列でエラーが発生しました: %v", err)
	}
	if postContent.String() != "" {
		t.Errorf("空文字列の場合のString() = %v, want %v", postContent.String(), "")
	}
}

func TestPostContent_WithWhitespace(t *testing.T) {
	// 空白文字のみの場合
	whitespaceContent := "   \n\t\r   "
	postContent, err := NewPostContent(whitespaceContent)
	if err != nil {
		t.Errorf("空白文字でエラーが発生しました: %v", err)
	}
	if postContent.String() != whitespaceContent {
		t.Errorf("空白文字の場合のString() = %v, want %v", postContent.String(), whitespaceContent)
	}
}