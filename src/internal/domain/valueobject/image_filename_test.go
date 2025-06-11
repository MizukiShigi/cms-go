package valueobject

import (
	"strings"
	"testing"
)

func TestNewImageFilename(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "正常ケース: JPGファイル",
			filename: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "正常ケース: JPEGファイル",
			filename: "photo.jpeg",
			wantErr:  false,
		},
		{
			name:     "正常ケース: PNGファイル",
			filename: "screenshot.png",
			wantErr:  false,
		},
		{
			name:     "正常ケース: GIFファイル",
			filename: "animation.gif",
			wantErr:  false,
		},
		{
			name:     "正常ケース: WEBPファイル",
			filename: "modern.webp",
			wantErr:  false,
		},
		{
			name:     "正常ケース: 大文字拡張子",
			filename: "IMAGE.JPG",
			wantErr:  false,
		},
		{
			name:     "正常ケース: 日本語ファイル名",
			filename: "画像.jpg",
			wantErr:  false,
		},
		{
			name:     "正常ケース: 数字を含む",
			filename: "img_001.png",
			wantErr:  false,
		},
		{
			name:     "正常ケース: 最大長（255文字）",
			filename: strings.Repeat("a", 251) + ".jpg",
			wantErr:  false,
		},
		{
			name:        "異常ケース: 空文字",
			filename:    "",
			wantErr:     true,
			expectedErr: "filename length invalid",
		},
		{
			name:        "異常ケース: 長すぎるファイル名（256文字）",
			filename:    strings.Repeat("a", 252) + ".jpg",
			wantErr:     true,
			expectedErr: "filename length invalid",
		},
		{
			name:        "異常ケース: 禁止文字を含む（スラッシュ）",
			filename:    "path/to/image.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（バックスラッシュ）",
			filename:    "path\\to\\image.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（コロン）",
			filename:    "C:image.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（アスタリスク）",
			filename:    "image*.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（クエスチョン）",
			filename:    "image?.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（ダブルクォート）",
			filename:    "\"image\".jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（不等号）",
			filename:    "<image>.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: 禁止文字を含む（パイプ）",
			filename:    "image|copy.jpg",
			wantErr:     true,
			expectedErr: "filename contains invalid characters",
		},
		{
			name:        "異常ケース: サポートされていない拡張子（txt）",
			filename:    "document.txt",
			wantErr:     true,
			expectedErr: "unsupported file extension",
		},
		{
			name:        "異常ケース: サポートされていない拡張子（pdf）",
			filename:    "file.pdf",
			wantErr:     true,
			expectedErr: "unsupported file extension",
		},
		{
			name:        "異常ケース: 拡張子なし",
			filename:    "filename",
			wantErr:     true,
			expectedErr: "unsupported file extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := NewImageFilename(tt.filename)

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

			if filename.String() != tt.filename {
				t.Errorf("ImageFilename.String() = %v, want %v", filename.String(), tt.filename)
			}
		})
	}
}

func TestImageFilename_String(t *testing.T) {
	filename := "test.jpg"
	imageFilename, err := NewImageFilename(filename)
	if err != nil {
		t.Fatalf("ImageFilename作成に失敗: %v", err)
	}

	if imageFilename.String() != filename {
		t.Errorf("String() = %v, want %v", imageFilename.String(), filename)
	}
}

func TestImageFilename_Equals(t *testing.T) {
	filename1, _ := NewImageFilename("image1.jpg")
	filename2, _ := NewImageFilename("image1.jpg")
	filename3, _ := NewImageFilename("image2.png")

	// 同じファイル名
	if !filename1.Equals(filename2) {
		t.Error("同じファイル名同士の比較がfalseになりました")
	}

	// 異なるファイル名
	if filename1.Equals(filename3) {
		t.Error("異なるファイル名同士の比較がtrueになりました")
	}

	// 自身との比較
	if !filename1.Equals(filename1) {
		t.Error("自身との比較がfalseになりました")
	}
}

func TestImageFilename_CaseInsensitiveExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "小文字拡張子",
			filename: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "大文字拡張子",
			filename: "image.JPG",
			wantErr:  false,
		},
		{
			name:     "混在拡張子",
			filename: "image.JpG",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewImageFilename(tt.filename)

			if tt.wantErr && err == nil {
				t.Error("エラーが期待されましたが、エラーが発生しませんでした")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
		})
	}
}
