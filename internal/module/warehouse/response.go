package warehouse

import (
	"github.com/google/uuid"
)

type WarehouseDetailResponse struct {
	ID       uuid.UUID `json:"ID"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Province string    `json:"province"`
	City     string    `json:"city"`
	ZipCode  string    `json:"zipCode"`
}
