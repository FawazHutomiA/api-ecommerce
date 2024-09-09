package repository

import (
	campaignEntity "example/internal/entity"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindAll() ([]campaignEntity.Campaign, error)
	FindByUserID(userID int) ([]campaignEntity.Campaign, error)
	FindByID(ID int) (campaignEntity.Campaign, error)
	Save(campaign campaignEntity.Campaign) (campaignEntity.Campaign, error)
	Update(campaign campaignEntity.Campaign) (campaignEntity.Campaign, error)
	CreateImage(campaignImage campaignEntity.CampaignImage) (campaignEntity.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func CampaignNewRepository(db *gorm.DB) *campaignRepository {
	return &campaignRepository{db}
}

func (r *campaignRepository) FindAll() ([]campaignEntity.Campaign, error) {
	var campaigns []campaignEntity.Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepository) FindByUserID(userID int) ([]campaignEntity.Campaign, error) {
	var campaigns []campaignEntity.Campaign

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *campaignRepository) FindByID(ID int) (campaignEntity.Campaign, error) {
	var campaign campaignEntity.Campaign

	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ? ", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepository) Save(campaign campaignEntity.Campaign) (campaignEntity.Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepository) Update(campaign campaignEntity.Campaign) (campaignEntity.Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *campaignRepository) CreateImage(campaignImage campaignEntity.CampaignImage) (campaignEntity.CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *campaignRepository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	// UPDATE campaign_images SET is_primary = false where campaign_id ? =

	err := r.db.Model(&campaignEntity.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
