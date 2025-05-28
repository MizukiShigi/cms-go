package valueobject

import (
	"strings"
	"testing"
)

func TestNewPostTitle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 通常のタイトル",
			title:   "これは通常のタイトルです",
			wantErr: false,
		},
		{
			name:    "正常ケース: 空のタイトル",
			title:   "",
			wantErr: false,
		},
		{
			name:    "正常ケース: 前後に空白があるタイトル",
			title:   "  タイトル  ",
			wantErr: false,
		},
		{
			name:    "正常ケース: 数字を含むタイトル",
			title:   "記事No.123",
			wantErr: false,
		},
		{
			name:    "正常ケース: 最大長のタイトル（200文字）",
			title:   strings.Repeat("あ", 200),
			wantErr: false,
		},
		{
			name:    "正常ケース: 英数字のタイトル",
			title:   "My Blog Post 2023",
			wantErr: false,
		},
		{
			name:    "正常ケース: 特殊文字を含むタイトル（許可されたもの）",
			title:   "ブログ記事：重要なお知らせ！",
			wantErr: false,
		},
		{
			name:        "異常ケース: 長すぎるタイトル（201文字）",
			title:       strings.Repeat("あ", 201),
			wantErr:     true,
			expectedErr: "Title is too long",
		},
		{
			name:        "異常ケース: 非常に長いタイトル",
			title:       strings.Repeat("テスト", 100),
			wantErr:     true,
			expectedErr: "Title is too long",
		},
		{
			name:        "異常ケース: 改行文字を含む",
			title:       "タイトル\n改行",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\\n'",
		},
		{
			name:        "異常ケース: タブ文字を含む",
			title:       "タイトル\tタブ",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\\t'",
		},
		{
			name:        "異常ケース: キャリッジリターンを含む",
			title:       "タイトル\r復帰",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\\r'",
		},
		{
			name:        "異常ケース: HTMLタグを含む",
			title:       "タイトル<script>alert('xss')</script>",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '<'",
		},
		{
			name:        "異常ケース: 引用符を含む",
			title:       "タイトル\"引用符",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\"'",
		},
		{
			name:        "異常ケース: 単一引用符を含む",
			title:       "タイトル'単一引用符",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\\''",
		},
		{
			name:        "異常ケース: バックスラッシュを含む",
			title:       "タイトル\\バックスラッシュ",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '\\\\'",
		},
		{
			name:        "異常ケース: スラッシュを含む",
			title:       "タイトル/スラッシュ",
			wantErr:     true,
			expectedErr: "Title contains forbidden character: '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewPostTitle(tt.title)

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

			// 前後の空白が削除されていることを確認
			expectedTitle := strings.TrimSpace(tt.title)
			if title.Value() != expectedTitle {
				t.Errorf("PostTitle.Value() = %v, want %v", title.Value(), expectedTitle)
			}
		})
	}
}

func TestPostTitle_Value(t *testing.T) {
	title := "テストタイトル"
	postTitle, err := NewPostTitle(title)
	if err != nil {
		t.Fatalf("PostTitle作成に失敗: %v", err)
	}

	if postTitle.Value() != title {
		t.Errorf("Value() = %v, want %v", postTitle.Value(), title)
	}
}

func TestPostTitle_String(t *testing.T) {
	title := "テストタイトル"
	postTitle, err := NewPostTitle(title)
	if err != nil {
		t.Fatalf("PostTitle作成に失敗: %v", err)
	}

	// HTMLエスケープされた形で格納されている
	if postTitle.String() == "" {
		t.Error("String()が空文字を返しました")
	}
}

func TestPostTitle_Equals(t *testing.T) {
	title1, _ := NewPostTitle("同じタイトル")
	title2, _ := NewPostTitle("同じタイトル")
	title3, _ := NewPostTitle("異なるタイトル")

	// 同じタイトル
	if !title1.Equals(title2) {
		t.Error("同じタイトル同士の比較がfalseになりました")
	}

	// 異なるタイトル
	if title1.Equals(title3) {
		t.Error("異なるタイトル同士の比較がtrueになりました")
	}

	// 自身との比較
	if !title1.Equals(title1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestPostTitle_TrimWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "前後に空白がある場合",
			input:    "  タイトル  ",
			expected: "タイトル",
		},
		{
			name:     "前に空白がある場合",
			input:    "  タイトル",
			expected: "タイトル",
		},
		{
			name:     "後に空白がある場合",
			input:    "タイトル  ",
			expected: "タイトル",
		},
		{
			name:     "空白なしの場合",
			input:    "タイトル",
			expected: "タイトル",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewPostTitle(tt.input)
			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			if title.Value() != tt.expected {
				t.Errorf("Value() = %v, want %v", title.Value(), tt.expected)
			}
		})
	}
}

func TestPostTitle_Boundary(t *testing.T) {
	// 境界値テスト：200文字ちょうど
	title200 := strings.Repeat("あ", 200)
	postTitle, err := NewPostTitle(title200)
	if err != nil {
		t.Errorf("200文字のタイトルでエラーが発生しました: %v", err)
	}
	if len(postTitle.Value()) != 200 {
		t.Errorf("200文字のタイトル長 = %d, want 200", len(postTitle.Value()))
	}

	// 境界値テスト：201文字
	title201 := strings.Repeat("あ", 201)
	_, err = NewPostTitle(title201)
	if err == nil {
		t.Error("201文字のタイトルでエラーが発生しませんでした")
	}
}

func TestPostTitle_HTMLEscaping(t *testing.T) {
	// HTMLエスケープのテスト
	title := "Title & Content"
	postTitle, err := NewPostTitle(title)
	if err != nil {
		t.Fatalf("PostTitle作成に失敗: %v", err)
	}

	// Value()は元の値を返す
	if postTitle.Value() != title {
		t.Errorf("Value() = %v, want %v", postTitle.Value(), title)
	}

	// String()はエスケープされた値を返す
	if postTitle.String() == title {
		t.Error("String()がエスケープされていません")
	}
}