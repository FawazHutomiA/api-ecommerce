package cloudinary

import (
	"example/pkg/app"
	"example/pkg/response"
	"example/pkg/validator"
	"net/http"
)

type CloudinaryHandler struct {
	App               app.AppConfig
	CloudinaryService CloudinaryService
}

func NewCloudinaryHandler(app app.AppConfig, CloudinaryService CloudinaryService) *CloudinaryHandler {
	return &CloudinaryHandler{
		App:               app,
		CloudinaryService: CloudinaryService,
	}
}

func (handler *CloudinaryHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Init
	var req uploadImageRequest
	var resp response.Response
	ctx := r.Context()

	r.ParseMultipartForm(10 << 20)
	formData := r.MultipartForm

	image := formData.File["image"]
	if len(image) > 0 {
		req.Image = image[0]
	}

	resp, errV := validator.ValidateFormRequest(r, &req)
	if errV != nil {
		resp.JSON(w)
		return
	}

	service, err := handler.CloudinaryService.UploadImage(ctx, req)
	if err.Errors != nil {
		handler.App.Logger.Error(err)
		resp = response.Error(err.Status, err.Message, err.Errors)
		resp.JSON(w)
		return
	}

	resp = response.Success(response.StatusOK, "Success", service)
	resp.JSON(w)
}
