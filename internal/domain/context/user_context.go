package context

import (
	"context"

	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
)

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserID).(string)
	if !ok {
		return "", myerror.NewMyError(myerror.UnauthorizedCode, "ユーザーが見つかりません")
	}

	return userID, nil
}
