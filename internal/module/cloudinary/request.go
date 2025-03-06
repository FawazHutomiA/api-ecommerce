package cloudinary

import "mime/multipart"

type uploadImageRequest struct {
	Image      *multipart.FileHeader `schema:"image" json:"image" validate:"omitempty,maxSizeFile=10,typeFile=ico png jpeg jpg webp svg"`
	ModuleName string                `schema:"moduleName" json:"moduleName"`
}
