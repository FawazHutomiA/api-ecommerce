package repository

import (
	transactionEntity "example/internal/entity"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepository interface {
	GetByCampaignID(campaignID int) ([]transactionEntity.Transaction, error)
	GetByUserID(userID int) ([]transactionEntity.Transaction, error)
	GetByID(ID int) (transactionEntity.Transaction, error)
	Save(transaction transactionEntity.Transaction) (transactionEntity.Transaction, error)
	Update(transaction transactionEntity.Transaction) (transactionEntity.Transaction, error)
}

func TransactionNewRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetByCampaignID(campaignID int) ([]transactionEntity.Transaction, error) {
	var transactions []transactionEntity.Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetByUserID(userID int) ([]transactionEntity.Transaction, error) {
	var transactions []transactionEntity.Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = true").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetByID(ID int) (transactionEntity.Transaction, error) {
	var transaction transactionEntity.Transaction

	err := r.db.Where("id = ? ", ID).Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) Save(transaction transactionEntity.Transaction) (transactionEntity.Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) Update(transaction transactionEntity.Transaction) (transactionEntity.Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
