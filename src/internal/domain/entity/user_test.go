package entity

import (
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "正常ケース: 有効なユーザー作成",
			userName: "testuser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:        "異常ケース: 空の名前",
			userName:    "",
			email:       "test@example.com",
			password:    "password123",
			wantErr:     true,
			expectedErr: "Name cannot be empty",
		},
		{
			name:        "異常ケース: 短いパスワード",
			userName:    "testuser",
			email:       "test@example.com",
			password:    "1234567",
			wantErr:     true,
			expectedErr: "Password must be at least 8 characters",
		},
		{
			name:        "異常ケース: 不正なメールアドレス",
			userName:    "testuser",
			email:       "invalid-email",
			password:    "password123",
			wantErr:     true,
			expectedErr: "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, emailErr := valueobject.NewEmail(tt.email)
			if emailErr != nil && !tt.wantErr {
				t.Errorf("メールアドレス作成でエラー: %v", emailErr)
				return
			}
			if emailErr != nil && tt.wantErr && tt.expectedErr == "Invalid email format" {
				// メールアドレス形式エラーの場合はここで期待通り
				return
			}

			user, err := NewUser(tt.userName, email, tt.password)

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

			if user == nil {
				t.Error("ユーザーがnilです")
				return
			}

			// フィールドの検証
			if user.Name != tt.userName {
				t.Errorf("Name = %v, want %v", user.Name, tt.userName)
			}

			if user.Email.String() != tt.email {
				t.Errorf("Email = %v, want %v", user.Email.String(), tt.email)
			}

			// パスワードがハッシュ化されていることを確認
			if user.Password == tt.password {
				t.Error("パスワードがハッシュ化されていません")
			}

			// IDが生成されていることを確認
			if user.ID.String() == "" {
				t.Error("ユーザーIDが生成されていません")
			}

			// タイムスタンプが設定されていることを確認
			if user.CreatedAt.IsZero() {
				t.Error("CreatedAtが設定されていません")
			}

			if user.UpdatedAt.IsZero() {
				t.Error("UpdatedAtが設定されていません")
			}
		})
	}
}

func TestUser_Authenticate(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")
	password := "password123"
	user, err := NewUser("testuser", email, password)
	if err != nil {
		t.Fatalf("ユーザー作成に失敗: %v", err)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "正常ケース: 正しいパスワード",
			password: "password123",
			want:     true,
		},
		{
			name:     "異常ケース: 間違ったパスワード",
			password: "wrongpassword",
			want:     false,
		},
		{
			name:     "異常ケース: 空のパスワード",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := user.Authenticate(tt.password)
			if got != tt.want {
				t.Errorf("Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}