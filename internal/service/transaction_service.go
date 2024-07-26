package service

import (
	entity "example/internal/entity"
	input "example/internal/input"
	repository "example/internal/repository"

	"errors"
	"strconv"
)

type transactionService struct {
	repository         repository.TransactionRepository
	campaignRepository repository.CampaignRepository
	paymentService     PaymentService
}

func TransactionNewService(repository repository.TransactionRepository, campaignRepository repository.CampaignRepository, paymentService PaymentService) *transactionService {
	return &transactionService{repository, campaignRepository, paymentService}
}

type TransactionService interface {
	GetTransactionsByCampaignID(input input.GetCampaignTransactionsInput) ([]entity.Transaction, error)
	GetTransactionsByUserID(userID int) ([]entity.Transaction, error)
	CreateTransaction(input input.CreateTransactionInput) (entity.Transaction, error)
	ProcessPayment(input input.TransactionNotificationInput) error
}

func (s *transactionService) GetTransactionsByCampaignID(input input.GetCampaignTransactionsInput) ([]entity.Transaction, error) {

	campaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []entity.Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []entity.Transaction{}, errors.New("Not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *transactionService) CreateTransaction(input input.CreateTransactionInput) (entity.Transaction, error) {
	transaction := entity.Transaction{}

	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := entity.PaymentTransaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentUrl, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentUrl

	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil

}

func (s *transactionService) ProcessPayment(input input.TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)

	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)

		if err != nil {
			return err
		}
	}

	return nil

}
