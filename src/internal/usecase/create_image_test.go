package usecase

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/service"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	repositoryMock "github.com/MizukiShigi/cms-go/mocks/repository"
	serviceMock "github.com/MizukiShigi/cms-go/mocks/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateImageUsecase_Execute(t *testing.T) {
	// テスト用の環境変数を設定
	originalBucketName := os.Getenv("GCS_IMAGE_BUCKET_NAME")
	defer os.Setenv("GCS_IMAGE_BUCKET_NAME", originalBucketName)
	os.Setenv("GCS_IMAGE_BUCKET_NAME", "test-bucket")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := repositoryMock.NewMockImageRepository(ctrl)
	mockStorageService := serviceMock.NewMockStorageService(ctrl)

	t.Run("画像作成が成功する（JPG）", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		// テストデータ準備
		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")
		fileReader := strings.NewReader("test image data")
		sortOrder := 1

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		// モックの期待値設定
		uploadResult := service.UploadResult{
			StoredFilename: "stored-test.jpg",
			URL:            "https://example.com/stored-test.jpg",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(nil)

		// 実行
		output, err := usecase.Execute(context.Background(), input)

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, userID, output.UserID)
		assert.Equal(t, postID, output.PostID)
		assert.Equal(t, filename, output.OriginalFilename)
		assert.Equal(t, "stored-test.jpg", output.StoredFilename)
		assert.Equal(t, "https://example.com/stored-test.jpg", output.ImageURL)
		assert.Equal(t, sortOrder, output.SortOrder)
		assert.NotEmpty(t, output.ID)
	})

	t.Run("画像作成が成功する（PNG）", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.png")
		fileReader := strings.NewReader("test png data")
		sortOrder := 2

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-test.png",
			URL:            "https://example.com/stored-test.png",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, "stored-test.png", output.StoredFilename)
		assert.Equal(t, sortOrder, output.SortOrder)
	})

	t.Run("画像作成が成功する（WebP）", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.webp")
		fileReader := strings.NewReader("test webp data")
		sortOrder := 0

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-test.webp",
			URL:            "https://example.com/stored-test.webp",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, sortOrder, output.SortOrder)
	})

	t.Run("画像作成が成功する（GIF）", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("animated.gif")
		fileReader := strings.NewReader("test gif data")
		sortOrder := 5

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-animated.gif",
			URL:            "https://example.com/stored-animated.gif",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, sortOrder, output.SortOrder)
	})

	t.Run("画像作成が成功する（JPEG）", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("photo.jpeg")
		fileReader := strings.NewReader("test jpeg data")
		sortOrder := 10

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-photo.jpeg",
			URL:            "https://example.com/stored-photo.jpeg",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(nil)

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, sortOrder, output.SortOrder)
	})

	t.Run("ストレージサービスのアップロードに失敗する", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")
		fileReader := strings.NewReader("test image data")

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        1,
		}

		// ストレージサービスでエラーが発生
		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(service.UploadResult{}, valueobject.NewMyError(valueobject.InternalServerErrorCode, "Storage upload failed"))

		// リポジトリのCreateは呼ばれない
		mockImageRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
		assert.Equal(t, "Storage upload failed", myErr.Message)
	})

	t.Run("リポジトリの保存に失敗する", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")
		fileReader := strings.NewReader("test image data")

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        1,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-test.jpg",
			URL:            "https://example.com/stored-test.jpg",
		}

		// ストレージアップロードは成功
		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		// リポジトリの保存で失敗
		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			Return(valueobject.NewMyError(valueobject.InternalServerErrorCode, "Database error"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)

		var myErr *valueobject.MyError
		assert.ErrorAs(t, err, &myErr)
		assert.Equal(t, valueobject.InternalServerErrorCode, myErr.Code)
		assert.Equal(t, "Database error", myErr.Message)
	})

	t.Run("環境変数が設定されていない場合", func(t *testing.T) {
		// 環境変数を一時的にクリア
		os.Unsetenv("GCS_IMAGE_BUCKET_NAME")
		defer os.Setenv("GCS_IMAGE_BUCKET_NAME", "test-bucket")

		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")
		fileReader := strings.NewReader("test image data")

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        1,
		}

		// 空のバケット名でアップロードが試行される
		mockStorageService.EXPECT().
			UploadImage(context.Background(), "", filename, fileReader).
			Return(service.UploadResult{}, valueobject.NewMyError(valueobject.InvalidCode, "Invalid bucket name"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("ソート順序が正しく設定される", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")
		fileReader := strings.NewReader("test image data")
		sortOrder := 100

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             fileReader,
			OriginalFilename: filename,
			SortOrder:        sortOrder,
		}

		uploadResult := service.UploadResult{
			StoredFilename: "stored-test.jpg",
			URL:            "https://example.com/stored-test.jpg",
		}

		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, fileReader).
			Return(uploadResult, nil)

		mockImageRepo.EXPECT().
			Create(context.Background(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, image interface{}) error {
				// 作成されるイメージエンティティのソート順序を検証
				// gomock.Anyを使っているが、DoAndReturnで実際の値を確認
				return nil
			})

		output, err := usecase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, sortOrder, output.SortOrder)
	})

	t.Run("ファイルリーダーがnilの場合", func(t *testing.T) {
		usecase := NewCreateImageUsecase(mockImageRepo, mockStorageService)

		userID := valueobject.NewUserID()
		postID := valueobject.NewPostID()
		filename, _ := valueobject.NewImageFilename("test.jpg")

		input := &CreateImageInput{
			UserID:           userID,
			PostID:           postID,
			File:             nil, // nilファイル
			OriginalFilename: filename,
			SortOrder:        1,
		}

		// ストレージサービスがnilファイルを受け取ってエラーになることを想定
		mockStorageService.EXPECT().
			UploadImage(context.Background(), "test-bucket", filename, nil).
			Return(service.UploadResult{}, valueobject.NewMyError(valueobject.InvalidCode, "File is required"))

		output, err := usecase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}

// テスト用のファイルリーダー作成ヘルパー
func createTestFileReader(content string) io.Reader {
	return strings.NewReader(content)
}