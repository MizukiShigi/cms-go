package valueobject

import (
	"strings"
	"testing"
)

func TestNewTagName(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 小文字の英字のみ",
			tag:     "technology",
			wantErr: false,
		},
		{
			name:    "正常ケース: 数字を含む",
			tag:     "tech2023",
			wantErr: false,
		},
		{
			name:    "正常ケース: ハイフンを含む",
			tag:     "web-development",
			wantErr: false,
		},
		{
			name:    "正常ケース: アンダースコアを含む",
			tag:     "go_lang",
			wantErr: false,
		},
		{
			name:    "正常ケース: 前後に空白があるタグ",
			tag:     "  programming  ",
			wantErr: false,
		},
		{
			name:    "正常ケース: 最大長のタグ（50文字）",
			tag:     strings.Repeat("a", 50),
			wantErr: false,
		},
		{
			name:    "正常ケース: 数字のみ",
			tag:     "2023",
			wantErr: false,
		},
		{
			name:        "正常ケース: 日本語を含む",
			tag:         "プログラミング",
			wantErr:     false,
		},
		{
			name:    "正常ケース: ハイフンとアンダースコアの組み合わせ",
			tag:     "web-dev_2023",
			wantErr: false,
		},
		{
			name:        "異常ケース: 長すぎるタグ（51文字）",
			tag:         strings.Repeat("a", 51),
			wantErr:     true,
			expectedErr: "TagName is too long",
		},
		{
			name:        "異常ケース: 非常に長いタグ",
			tag:         strings.Repeat("tag", 25),
			wantErr:     true,
			expectedErr: "TagName is too long",
		},
		{
			name:        "異常ケース: スペースを含む",
			tag:         "web development",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
		{
			name:        "異常ケース: 特殊文字を含む",
			tag:         "tech!",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
		{
			name:        "異常ケース: ドットを含む",
			tag:         "node.js",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
		{
			name:        "異常ケース: カンマを含む",
			tag:         "web,dev",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
		{
			name:        "異常ケース: アットマークを含む",
			tag:         "@tag",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
		{
			name:        "異常ケース: スラッシュを含む",
			tag:         "web/dev",
			wantErr:     true,
			expectedErr: "TagName can only contain Japanese characters, letters, numbers, hyphens, and underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagName, err := NewTagName(tt.tag)

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

			// 正規化された値が返されることを確認
			expectedNormalized := strings.ToLower(strings.TrimSpace(tt.tag))
			if tagName.String() != expectedNormalized {
				t.Errorf("TagName.String() = %v, want %v", tagName.String(), expectedNormalized)
			}
		})
	}
}

func TestTagName_String(t *testing.T) {
	tag := "programming"
	tagName, err := NewTagName(tag)
	if err != nil {
		t.Fatalf("TagName作成に失敗: %v", err)
	}

	if tagName.String() != tag {
		t.Errorf("String() = %v, want %v", tagName.String(), tag)
	}
}

func TestTagName_Equals(t *testing.T) {
	tag1, _ := NewTagName("golang")
	tag2, _ := NewTagName("golang")
	tag3, _ := NewTagName("javascript")

	// 同じタグ名
	if !tag1.Equals(tag2) {
		t.Error("同じタグ名同士の比較がfalseになりました")
	}

	// 異なるタグ名
	if tag1.Equals(tag3) {
		t.Error("異なるタグ名同士の比較がtrueになりました")
	}

	// 自身との比較
	if !tag1.Equals(tag1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestTagName_Normalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "前後の空白が削除される",
			input:    "  golang  ",
			expected: "golang",
		},
		{
			name:     "すでに正規化されている場合",
			input:    "python",
			expected: "python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagName, err := NewTagName(tt.input)
			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if tagName.String() != tt.expected {
				t.Errorf("String() = %v, want %v", tagName.String(), tt.expected)
			}
		})
	}
}

func TestTagName_Boundary(t *testing.T) {
	// 境界値テスト：50文字ちょうど
	tag50 := strings.Repeat("a", 50)
	tagName, err := NewTagName(tag50)
	if err != nil {
		t.Errorf("50文字のタグでエラーが発生しました: %v", err)
	}
	if len(tagName.String()) != 50 {
		t.Errorf("50文字のタグ長 = %d, want 50", len(tagName.String()))
	}

	// 境界値テスト：51文字
	tag51 := strings.Repeat("a", 51)
	_, err = NewTagName(tag51)
	if err == nil {
		t.Error("51文字のタグでエラーが発生しませんでした")
	}
}

func TestTagName_ValidCharacters(t *testing.T) {
	validChars := []string{
		"a", "z", "0", "9", "-", "_",
		"abc123", "test-tag", "tag_name", "a1b2c3",
	}

	for _, char := range validChars {
		t.Run("Valid: "+char, func(t *testing.T) {
			_, err := NewTagName(char)
			if err != nil {
				t.Errorf("有効な文字列 %s でエラーが発生しました: %v", char, err)
			}
		})
	}
}

func TestTagName_InvalidCharacters(t *testing.T) {
	invalidChars := []string{
		" ", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")",
		"+", "=", "[", "]", "{", "}", "|", "\\", ":", ";", "\"",
		"'", "<", ">", ",", ".", "?", "/",
	}

	for _, char := range invalidChars {
		t.Run("Invalid: "+char, func(t *testing.T) {
			_, err := NewTagName("test" + char + "tag")
			if err == nil {
				t.Errorf("無効な文字 %s を含むタグでエラーが発生しませんでした", char)
			}
		})
	}
}

func TestTagName_EmptyString(t *testing.T) {
	// 空文字列の場合（正規化後も空文字列）
	_, err := NewTagName("")
	if err == nil {
		t.Errorf("空文字列でエラーが発生しました: %v", err)
	}
}

func TestTagName_OnlyWhitespace(t *testing.T) {
	// 空白のみの場合（正規化後は空文字列）
	_, err := NewTagName("   ")
	if err == nil {
		t.Errorf("空白のみでエラーが発生しませんでした: %v", err)
	}
}
