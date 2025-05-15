package context

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserID).(string)
	if !ok {
		return "", valueobject.NewMyError(valueobject.UnauthorizedCode, "ユーザーが見つかりません")
	}

	return userID, nil
}
