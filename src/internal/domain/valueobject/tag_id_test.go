package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTagID(t *testing.T) {
	tagID := NewTagID()

	// 空でないことを確認
	if tagID.String() == "" {
		t.Error("TagIDが空文字です")
	}

	// UUID形式であることを確認
	_, err := uuid.Parse(tagID.String())
	if err != nil {
		t.Errorf("TagIDが有効なUUID形式ではありません: %v", err)
	}

	// 2回連続で生成した際に異なる値であることを確認
	tagID2 := NewTagID()
	if tagID.Equals(tagID2) {
		t.Error("2回連続で生成したTagIDが同じ値になりました")
	}
}

func TestParseTagID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
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
			name:    "異常ケース: 無効なUUID形式",
			input:   "invalid-uuid",
			wantErr: true,
		},
		{
			name:    "異常ケース: 空文字",
			input:   "",
			wantErr: true,
		},
		{
			name:    "異常ケース: 短いUUID",
			input:   "550e8400-e29b-41d4-a716",
			wantErr: true,
		},
		{
			name:    "異常ケース: 長いUUID",
			input:   "550e8400-e29b-41d4-a716-446655440000-extra",
			wantErr: true,
		},
		{
			name:    "異常ケース: 大文字小文字混在（無効な文字）",
			input:   "550e8400-e29b-41d4-a716-44665544000G",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagID, err := ParseTagID(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// 正規化されたUUID形式であることを確認
			expectedNormalized := "550e8400-e29b-41d4-a716-446655440000"
			if tt.input == "550e8400-e29b-41d4-a716-446655440000" && tagID.String() != expectedNormalized {
				t.Errorf("ParseTagID() = %v, want %v", tagID.String(), expectedNormalized)
			}
		})
	}
}

func TestTagID_String(t *testing.T) {
	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	tagID, err := ParseTagID(expectedUUID)
	if err != nil {
		t.Fatalf("TagID作成に失敗: %v", err)
	}

	if tagID.String() != expectedUUID {
		t.Errorf("String() = %v, want %v", tagID.String(), expectedUUID)
	}
}

func TestTagID_Equals(t *testing.T) {
	uuid1 := "550e8400-e29b-41d4-a716-446655440000"
	uuid2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	tagID1, _ := ParseTagID(uuid1)
	tagID2, _ := ParseTagID(uuid1) // 同じUUID
	tagID3, _ := ParseTagID(uuid2) // 異なるUUID

	// 同じUUID
	if !tagID1.Equals(tagID2) {
		t.Error("同じUUID同士の比較がfalseになりました")
	}

	// 異なるUUID
	if tagID1.Equals(tagID3) {
		t.Error("異なるUUID同士の比較がtrueになりました")
	}

	// 自身との比較
	if !tagID1.Equals(tagID1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestTagID_Generation(t *testing.T) {
	// 複数のTagIDを生成して重複がないことを確認
	tagIDs := make(map[string]bool)
	const numTests = 100

	for i := 0; i < numTests; i++ {
		tagID := NewTagID()
		idStr := tagID.String()

		if tagIDs[idStr] {
			t.Errorf("重複したTagIDが生成されました: %s", idStr)
		}
		tagIDs[idStr] = true

		// 各IDがUUID形式であることを確認
		_, err := uuid.Parse(idStr)
		if err != nil {
			t.Errorf("生成されたTagIDが有効なUUID形式ではありません: %s, error: %v", idStr, err)
		}
	}
}

func TestTagID_ParseAndEquals(t *testing.T) {
	// 生成されたTagIDをパースして同一性を確認
	originalID := NewTagID()
	parsedID, err := ParseTagID(originalID.String())

	if err != nil {
		t.Errorf("生成されたTagIDのパースに失敗: %v", err)
	}

	if !originalID.Equals(parsedID) {
		t.Error("生成されたTagIDとパースされたTagIDが一致しません")
	}
}

func TestTagID_EmptyString(t *testing.T) {
	// 空文字をパースした場合のテスト
	_, err := ParseTagID("")
	if err == nil {
		t.Error("空文字のパースでエラーが発生しませんでした")
	}
}

func TestTagID_ValidUUIDFormats(t *testing.T) {
	// 様々な有効なUUID形式をテスト
	validUUIDs := []string{
		"00000000-0000-0000-0000-000000000000", // nil UUID
		"123e4567-e89b-12d3-a456-426614174000", // version 1
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8", // version 1
		"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", // version 4
		"87b3042c-cc7e-4b1a-9a47-0a3f5fd2e2ad", // version 4
	}

	for _, uuidStr := range validUUIDs {
		t.Run("Valid UUID: "+uuidStr, func(t *testing.T) {
			tagID, err := ParseTagID(uuidStr)
			if err != nil {
				t.Errorf("有効なUUID %s のパースに失敗: %v", uuidStr, err)
			}
			if tagID.String() != uuidStr {
				t.Errorf("パース結果 = %v, want %v", tagID.String(), uuidStr)
			}
		})
	}
}
