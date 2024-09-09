package service

import (
	"example/internal/entity"

	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentService struct {
}

type PaymentService interface {
	GetPaymentURL(transaction entity.PaymentTransaction, user entity.User) (string, error)
}

func PaymentNewService() *paymentService {
	return &paymentService{}
}

func (s *paymentService) GetPaymentURL(transaction entity.PaymentTransaction, user entity.User) (string, error) {
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	// midtransClientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	midtrans.Environment = midtrans.Sandbox

	snapGateway := snap.Client{}
	snapGateway.New(midtransServerKey, midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenRespon, err := snapGateway.CreateTransaction(req)

	if err != nil {
		return "", err
	}

	return snapTokenRespon.RedirectURL, nil
}
