package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MizukiShigi/cms-go/infrastructure/db/sqlboiler/models"
	"github.com/MizukiShigi/cms-go/internal/domain/entity"
	"github.com/MizukiShigi/cms-go/internal/domain/myerror"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	dbUser := &models.User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email.String(),
		Password:  user.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := dbUser.Insert(ctx, ur.db, boil.Infer()); err != nil {
		// PostgreSQLの一意制約違反のエラーをチェック
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return myerror.NewMyError(myerror.ConflictCode, "User with this email already exists")
		}
		return myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to save user")
	}

	return nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	dbUser, err := models.Users(qm.Where(models.UserColumns.Email+" = ?", email.String())).One(ctx, ur.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myerror.NewMyError(myerror.NotFoundCode, "User not found")
		}
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to find user")
	}

	voUserID, err := valueobject.ParseUserID(dbUser.ID)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to parse user ID")
	}

	voEmail, err := valueobject.NewEmail(dbUser.Email)
	if err != nil {
		return nil, myerror.NewMyError(myerror.InternalServerErrorCode, "Failed to parse email")
	}

	return &entity.User{
		ID:        voUserID,
		Name:      dbUser.Name,
		Email:     voEmail,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}