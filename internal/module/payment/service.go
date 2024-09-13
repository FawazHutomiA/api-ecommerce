package payment

import (
	"example/internal/entity"

	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentService struct {
}

type PaymentService interface {
	GetPaymentURL(transaction entity.PaymentRequest, user entity.User) (string, error)
}

func PaymentNewService() *paymentService {
	return &paymentService{}
}

func (s *paymentService) GetPaymentURL(transaction entity.PaymentRequest, user entity.User) (string, error) {
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
			OrderID:  transaction.ID.String(),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenRespon, err := snapGateway.CreateTransaction(req)

	if err != nil {
		return "", err
	}

	return snapTokenRespon.RedirectURL, nil
}
