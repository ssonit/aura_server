package biz

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/media/models"
	"github.com/ssonit/aura_server/internal/media/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	cloudinaryCloudName = common.EnvConfig("CLOUDINARY_CLOUD_NAME", "")
	cloudinaryAPIKey    = common.EnvConfig("CLOUDINARY_API_KEY", "")
	cloudinaryAPISecret = common.EnvConfig("CLOUDINARY_API_SECRET", "")
)

type service struct {
	store utils.MediaStore
}

func NewService(store utils.MediaStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) GetMedia(ctx context.Context, id string) (*models.Media, error) {
	return s.store.GetMedia(ctx, id)
}

// UploadImage uploads an image to the server
func (s *service) UploadImage(ctx context.Context, file *multipart.FileHeader) (primitive.ObjectID, error) {
	f, err := file.Open()
	if err != nil {
		return primitive.NilObjectID, utils.ErrUnableToOpenFile
	}
	defer f.Close()

	cld, err := cloudinary.NewFromParams(
		cloudinaryCloudName,
		cloudinaryAPIKey,
		cloudinaryAPISecret,
	)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCldNewFromParams
	}

	publicID := common.GeneratePublicID()

	res, err := cld.Upload.Upload(ctx, f, uploader.UploadParams{
		PublicID:       publicID,
		Transformation: "f_auto",
	})

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotUploadCld
	}

	media := &models.MediaCreation{
		Url:       res.SecureURL,
		SecureUrl: res.SecureURL,
		PublicId:  res.PublicID,
		Format:    res.Format,
		Width:     res.Width,
		Height:    res.Height,
	}

	data, err := s.store.UploadImage(ctx, media)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil
}
