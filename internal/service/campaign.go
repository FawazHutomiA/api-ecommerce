package service

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"

	campaignEntity "example/internal/entity"
	campaignInput "example/internal/input"
	campaignRepository "example/internal/repository"
)

type CampaignService interface {
	GetCampaigns(userID int) ([]campaignEntity.Campaign, error)
	GetCampaignByID(input campaignInput.GetCampaignDetailInput) (campaignEntity.Campaign, error)
	CreateCampaign(input campaignInput.CreateCampaignInput) (campaignEntity.Campaign, error)
	UpdateCampaign(InputID campaignInput.GetCampaignDetailInput, InputData campaignInput.CreateCampaignInput) (campaignEntity.Campaign, error)
	SaveCampaignImage(input campaignInput.CreateCampaignImageInput, fileLocation string) (campaignEntity.CampaignImage, error)
}

type campaignService struct {
	repository campaignRepository.CampaignRepository
}

func CampaignNewService(repository campaignRepository.CampaignRepository) *campaignService {
	return &campaignService{repository}
}

func (s *campaignService) GetCampaigns(userID int) ([]campaignEntity.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *campaignService) GetCampaignByID(input campaignInput.GetCampaignDetailInput) (campaignEntity.Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignService) CreateCampaign(input campaignInput.CreateCampaignInput) (campaignEntity.Campaign, error) {
	campaign := campaignEntity.Campaign{}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.UserID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.UserID)

	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Save(campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *campaignService) UpdateCampaign(InputID campaignInput.GetCampaignDetailInput, InputData campaignInput.CreateCampaignInput) (campaignEntity.Campaign, error) {
	campaign, err := s.repository.FindByID(InputID.ID)

	if err != nil {
		return campaign, err
	}

	if campaign.UserID != InputData.UserID {
		return campaign, errors.New("Not an owner of the campaign")
	}

	campaign.Name = InputData.Name
	campaign.ShortDescription = InputData.ShortDescription
	campaign.Description = InputData.Description
	campaign.Perks = InputData.Perks
	campaign.GoalAmount = InputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)

	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *campaignService) SaveCampaignImage(input campaignInput.CreateCampaignImageInput, fileLocation string) (campaignEntity.CampaignImage, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)

	if err != nil {
		return campaignEntity.CampaignImage{}, err
	}

	if campaign.UserID != input.UserID {
		return campaignEntity.CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	if input.IsPrimary {
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)

		if err != nil {
			return campaignEntity.CampaignImage{}, err
		}
	}

	campaignImage := campaignEntity.CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = input.IsPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.CreateImage(campaignImage)

	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
