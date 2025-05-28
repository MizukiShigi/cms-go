package valueobject

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "正常ケース: 有効なメールアドレス",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "正常ケース: サブドメインを含むメールアドレス",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "正常ケース: 数字を含むメールアドレス",
			email:   "user123@example123.com",
			wantErr: false,
		},
		{
			name:    "正常ケース: ハイフンを含むメールアドレス",
			email:   "user-name@example-domain.com",
			wantErr: false,
		},
		{
			name:    "正常ケース: プラス記号を含むメールアドレス",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:        "異常ケース: @マークなし",
			email:       "testexample.com",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: ドメインなし",
			email:       "test@",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: ローカル部なし",
			email:       "@example.com",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: 拡張子なし",
			email:       "test@example",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: 空文字",
			email:       "",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: スペースを含む",
			email:       "test @example.com",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
		{
			name:        "異常ケース: 短すぎる拡張子",
			email:       "test@example.c",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)

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

			if email.String() != tt.email {
				t.Errorf("Email.String() = %v, want %v", email.String(), tt.email)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("メールアドレス作成に失敗: %v", err)
	}

	expected := "test@example.com"
	if email.String() != expected {
		t.Errorf("String() = %v, want %v", email.String(), expected)
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("different@example.com")

	// 同じメールアドレス
	if !email1.Equals(email2) {
		t.Error("同じメールアドレス同士の比較がfalseになりました")
	}

	// 異なるメールアドレス
	if email1.Equals(email3) {
		t.Error("異なるメールアドレス同士の比較がtrueになりました")
	}

	// 自身との比較
	if !email1.Equals(email1) {
		t.Error("自身との比較がfalseになりました")
	}
}
