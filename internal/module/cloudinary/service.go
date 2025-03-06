package cloudinary

import (
	"example/pkg/app"
	"example/pkg/cloudinary"
	"example/pkg/exception"
	"example/pkg/response"
	"example/pkg/strings"
	"fmt"
	"path/filepath"

	"context"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, params uploadImageRequest) (resp uploadImageResponse, errData exception.Error)
}

type cloudinaryService struct {
	app app.AppConfig
}

func NewCloudinaryService(app app.AppConfig) CloudinaryService {
	return &cloudinaryService{
		app: app,
	}
}

func (uc *cloudinaryService) UploadImage(ctx context.Context, params uploadImageRequest) (resp uploadImageResponse, errData exception.Error) {
	var image *string
	randomString := strings.RandomString(5)

	// Buka file dari request
	file, err := params.Image.Open()
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusInternalServerError,
			Message: "Failed to open Image",
			Errors:  exception.ErrInternalServer,
		}
	}
	defer file.Close()

	// Tentukan path di Cloudinary
	fileName := params.Image.Filename[:len(params.Image.Filename)-len(filepath.Ext(params.Image.Filename))]
	cloudinaryPath := fmt.Sprintf("%s/%s/%s", params.ModuleName, randomString, fileName)

	// Upload file ke Cloudinary (menggunakan io.Reader, bukan file path)
	url, err := cloudinary.UploadToCloudinary(file, cloudinaryPath)
	if err != nil {
		return resp, exception.Error{
			Status:  response.StatusBadRequest,
			Message: "Error upload",
			Errors:  exception.ErrBadRequest,
		}
	}

	image = &url

	resp = uploadImageResponse{
		Image: image,
	}

	return resp, errData
}
