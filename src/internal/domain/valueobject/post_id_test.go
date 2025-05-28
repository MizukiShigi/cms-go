package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewPostID(t *testing.T) {
	postID := NewPostID()

	// 空でないことを確認
	if postID.String() == "" {
		t.Error("PostIDが空文字です")
	}

	// UUID形式であることを確認
	_, err := uuid.Parse(postID.String())
	if err != nil {
		t.Errorf("PostIDが有効なUUID形式ではありません: %v", err)
	}

	// 2回連続で生成した際に異なる値であることを確認
	postID2 := NewPostID()
	if postID.Equals(postID2) {
		t.Error("2回連続で生成したPostIDが同じ値になりました")
	}
}

func TestParsePostID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 有効なUUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "正常ケース: 別の有効なUUID",
			input:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantErr: false,
		},
		{
			name:    "正常ケース: ランダムに生成されたUUID",
			input:   uuid.New().String(),
			wantErr: false,
		},
		{
			name:        "異常ケース: 無効なUUID形式",
			input:       "invalid-uuid",
			wantErr:     true,
			expectedErr: "Invalid post ID",
		},
		{
			name:        "異常ケース: 空文字",
			input:       "",
			wantErr:     true,
			expectedErr: "Invalid post ID",
		},
		{
			name:        "異常ケース: 短いUUID",
			input:       "550e8400-e29b-41d4-a716",
			wantErr:     true,
			expectedErr: "Invalid post ID",
		},
		{
			name:        "異常ケース: 長いUUID",
			input:       "550e8400-e29b-41d4-a716-446655440000-extra",
			wantErr:     true,
			expectedErr: "Invalid post ID",
		},
		{
			name:        "異常ケース: 大文字小文字混在（無効な文字）",
			input:       "550e8400-e29b-41d4-a716-44665544000G",
			wantErr:     true,
			expectedErr: "Invalid post ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postID, err := ParsePostID(tt.input)

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

			// 正規化されたUUID形式であることを確認
			expectedNormalized := "550e8400-e29b-41d4-a716-446655440000"
			if tt.input == "550e8400-e29b-41d4-a716-446655440000" && postID.String() != expectedNormalized {
				t.Errorf("ParsePostID() = %v, want %v", postID.String(), expectedNormalized)
			}
		})
	}
}

func TestPostID_String(t *testing.T) {
	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	postID, err := ParsePostID(expectedUUID)
	if err != nil {
		t.Fatalf("PostID作成に失敗: %v", err)
	}

	if postID.String() != expectedUUID {
		t.Errorf("String() = %v, want %v", postID.String(), expectedUUID)
	}
}

func TestPostID_Equals(t *testing.T) {
	uuid1 := "550e8400-e29b-41d4-a716-446655440000"
	uuid2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	postID1, _ := ParsePostID(uuid1)
	postID2, _ := ParsePostID(uuid1) // 同じUUID
	postID3, _ := ParsePostID(uuid2) // 異なるUUID

	// 同じUUID
	if !postID1.Equals(postID2) {
		t.Error("同じUUID同士の比較がfalseになりました")
	}

	// 異なるUUID
	if postID1.Equals(postID3) {
		t.Error("異なるUUID同士の比較がtrueになりました")
	}

	// 自身との比較
	if !postID1.Equals(postID1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestPostID_Generation(t *testing.T) {
	// 複数のPostIDを生成して重複がないことを確認
	postIDs := make(map[string]bool)
	const numTests = 100

	for i := 0; i < numTests; i++ {
		postID := NewPostID()
		idStr := postID.String()

		if postIDs[idStr] {
			t.Errorf("重複したPostIDが生成されました: %s", idStr)
		}
		postIDs[idStr] = true

		// 各IDがUUID形式であることを確認
		_, err := uuid.Parse(idStr)
		if err != nil {
			t.Errorf("生成されたPostIDが有効なUUID形式ではありません: %s, error: %v", idStr, err)
		}
	}
}

func TestPostID_ParseAndEquals(t *testing.T) {
	// 生成されたPostIDをパースして同一性を確認
	originalID := NewPostID()
	parsedID, err := ParsePostID(originalID.String())

	if err != nil {
		t.Errorf("生成されたPostIDのパースに失敗: %v", err)
	}

	if !originalID.Equals(parsedID) {
		t.Error("生成されたPostIDとパースされたPostIDが一致しません")
	}
}
