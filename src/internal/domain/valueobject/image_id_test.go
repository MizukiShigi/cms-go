package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewImageID(t *testing.T) {
	imageID := NewImageID()

	// 空でないことを確認
	if imageID.String() == "" {
		t.Error("ImageIDが空文字です")
	}

	// UUID形式であることを確認
	_, err := uuid.Parse(imageID.String())
	if err != nil {
		t.Errorf("ImageIDが有効なUUID形式ではありません: %v", err)
	}

	// 2回連続で生成した際に異なる値であることを確認
	imageID2 := NewImageID()
	if imageID.Equals(imageID2) {
		t.Error("2回連続で生成したImageIDが同じ値になりました")
	}
}

func TestParseImageID(t *testing.T) {
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
			expectedErr: "Invalid image ID",
		},
		{
			name:        "異常ケース: 空文字",
			input:       "",
			wantErr:     true,
			expectedErr: "Invalid image ID",
		},
		{
			name:        "異常ケース: 短いUUID",
			input:       "550e8400-e29b-41d4-a716",
			wantErr:     true,
			expectedErr: "Invalid image ID",
		},
		{
			name:        "異常ケース: 長いUUID",
			input:       "550e8400-e29b-41d4-a716-446655440000-extra",
			wantErr:     true,
			expectedErr: "Invalid image ID",
		},
		{
			name:        "異常ケース: 大文字小文字混在（無効な文字）",
			input:       "550e8400-e29b-41d4-a716-44665544000G",
			wantErr:     true,
			expectedErr: "Invalid image ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageID, err := ParseImageID(tt.input)

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
			if tt.input == "550e8400-e29b-41d4-a716-446655440000" && imageID.String() != expectedNormalized {
				t.Errorf("ParseImageID() = %v, want %v", imageID.String(), expectedNormalized)
			}
		})
	}
}

func TestImageID_String(t *testing.T) {
	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	imageID, err := ParseImageID(expectedUUID)
	if err != nil {
		t.Fatalf("ImageID作成に失敗: %v", err)
	}

	if imageID.String() != expectedUUID {
		t.Errorf("String() = %v, want %v", imageID.String(), expectedUUID)
	}
}

func TestImageID_Equals(t *testing.T) {
	uuid1 := "550e8400-e29b-41d4-a716-446655440000"
	uuid2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	imageID1, _ := ParseImageID(uuid1)
	imageID2, _ := ParseImageID(uuid1) // 同じUUID
	imageID3, _ := ParseImageID(uuid2) // 異なるUUID

	// 同じUUID
	if !imageID1.Equals(imageID2) {
		t.Error("同じUUID同士の比較がfalseになりました")
	}

	// 異なるUUID
	if imageID1.Equals(imageID3) {
		t.Error("異なるUUID同士の比較がtrueになりました")
	}

	// 自身との比較
	if !imageID1.Equals(imageID1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestImageID_Generation(t *testing.T) {
	// 複数のImageIDを生成して重複がないことを確認
	imageIDs := make(map[string]bool)
	const numTests = 100

	for i := 0; i < numTests; i++ {
		imageID := NewImageID()
		idStr := imageID.String()

		if imageIDs[idStr] {
			t.Errorf("重複したImageIDが生成されました: %s", idStr)
		}
		imageIDs[idStr] = true

		// 各IDがUUID形式であることを確認
		_, err := uuid.Parse(idStr)
		if err != nil {
			t.Errorf("生成されたImageIDが有効なUUID形式ではありません: %s, error: %v", idStr, err)
		}
	}
}

func TestImageID_ParseAndEquals(t *testing.T) {
	// 生成されたImageIDをパースして同一性を確認
	originalID := NewImageID()
	parsedID, err := ParseImageID(originalID.String())

	if err != nil {
		t.Errorf("生成されたImageIDのパースに失敗: %v", err)
	}

	if !originalID.Equals(parsedID) {
		t.Error("生成されたImageIDとパースされたImageIDが一致しません")
	}
}
