package cloudinary

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryConfig menyimpan konfigurasi Cloudinary
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// UploadToCloudinary mengupload file ke Cloudinary
func UploadToCloudinary(file multipart.File, fileName string) (string, error) {
	// Ganti dengan kredensial Cloudinary
	cld, err := cloudinary.NewFromParams("dmwmgaqyx", "979243457597517", "DQg8lW6cDhNL_36w3aTGPVcZKN8")
	if err != nil {
		log.Printf("Failed to initiate Cloudinary: %v", err)
		return "", err
	}

	// Upload file langsung dari io.Reader
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileName, // Nama file yang akan disimpan di Cloudinary
	})
	if err != nil {
		log.Printf("Failed to upload: %v", err)
		return "", err
	}

	return uploadResult.SecureURL, nil
}
