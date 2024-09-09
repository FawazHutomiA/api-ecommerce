package campaign

import (
	"example/internal/helper"

	campaignInput "example/internal/input"
	campaignService "example/internal/service"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CampaignHandler struct {
	service campaignService.CampaignService
}

func NewCampaignHandler(service campaignService.CampaignService) *CampaignHandler {
	return &CampaignHandler{service}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	// Query param
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input campaignInput.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input campaignInput.CreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	getUserID, _ := c.Get("userID")
	userID := getUserID.(int)
	input.UserID = userID

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaignInput.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaignInput.CreateCampaignInput

	err = c.ShouldBind(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	getUserID, _ := c.Get("userID")
	userID := getUserID.(int)
	inputData.UserID = userID

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) UploadImage(c *gin.Context) {
	var input campaignInput.CreateCampaignImageInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err.(validator.ValidationErrors))

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	getUserID, _ := c.Get("userID")
	userID := getUserID.(int)
	input.UserID = userID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
