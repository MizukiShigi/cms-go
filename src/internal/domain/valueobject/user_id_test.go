package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	userID := NewUserID()

	// 空でないことを確認
	if userID.String() == "" {
		t.Error("UserIDが空文字です")
	}

	// UUID形式であることを確認
	_, err := uuid.Parse(userID.String())
	if err != nil {
		t.Errorf("UserIDが有効なUUID形式ではありません: %v", err)
	}

	// 2回連続で生成した際に異なる値であることを確認
	userID2 := NewUserID()
	if userID.Equals(userID2) {
		t.Error("2回連続で生成したUserIDが同じ値になりました")
	}
}

func TestParseUserID(t *testing.T) {
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
			name:        "異常ケース: 無効なUUID形式",
			input:       "invalid-uuid",
			wantErr:     true,
			expectedErr: "Invalid user ID",
		},
		{
			name:        "異常ケース: 空文字",
			input:       "",
			wantErr:     true,
			expectedErr: "Invalid user ID",
		},
		{
			name:        "異常ケース: 短いUUID",
			input:       "550e8400-e29b-41d4-a716",
			wantErr:     true,
			expectedErr: "Invalid user ID",
		},
		{
			name:        "異常ケース: 長いUUID",
			input:       "550e8400-e29b-41d4-a716-446655440000-extra",
			wantErr:     true,
			expectedErr: "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ParseUserID(tt.input)

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
			if tt.input == "550e8400-e29b-41d4-a716-446655440000" && userID.String() != expectedNormalized {
				t.Errorf("ParseUserID() = %v, want %v", userID.String(), expectedNormalized)
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	userID, err := ParseUserID(expectedUUID)
	if err != nil {
		t.Fatalf("UserID作成に失敗: %v", err)
	}

	if userID.String() != expectedUUID {
		t.Errorf("String() = %v, want %v", userID.String(), expectedUUID)
	}
}

func TestUserID_Equals(t *testing.T) {
	uuid1 := "550e8400-e29b-41d4-a716-446655440000"
	uuid2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	userID1, _ := ParseUserID(uuid1)
	userID2, _ := ParseUserID(uuid1) // 同じUUID
	userID3, _ := ParseUserID(uuid2) // 異なるUUID

	// 同じUUID
	if !userID1.Equals(userID2) {
		t.Error("同じUUID同士の比較がfalseになりました")
	}

	// 異なるUUID
	if userID1.Equals(userID3) {
		t.Error("異なるUUID同士の比較がtrueになりました")
	}

	// 自身との比較
	if !userID1.Equals(userID1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestUserID_Generation(t *testing.T) {
	// 複数のUserIDを生成して重複がないことを確認
	userIDs := make(map[string]bool)
	const numTests = 100

	for i := 0; i < numTests; i++ {
		userID := NewUserID()
		idStr := userID.String()

		if userIDs[idStr] {
			t.Errorf("重複したUserIDが生成されました: %s", idStr)
		}
		userIDs[idStr] = true

		// 各IDがUUID形式であることを確認
		_, err := uuid.Parse(idStr)
		if err != nil {
			t.Errorf("生成されたUserIDが有効なUUID形式ではありません: %s, error: %v", idStr, err)
		}
	}
}
