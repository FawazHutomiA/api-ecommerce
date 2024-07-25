package payment

import (
	"example/module/user"
	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
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
