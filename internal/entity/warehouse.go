package entity

import "github.com/google/uuid"

type Warehouse struct {
	ID       uuid.UUID `db:"id" json:"ID"`
	Name     string    `db:"name" json:"name"`
	Address  string    `db:"address" json:"address"`
	Province string    `db:"province" json:"province"`
	City     string    `db:"city" json:"city"`
	ZipCode  string    `db:"zip_code" json:"zipCode"`
}

func (a *Warehouse) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.Name,
		a.Address,
		a.Province,
		a.City,
		a.ZipCode,
	}
}

func (a *Warehouse) ToUpdate() []interface{} {
	return []interface{}{
		a.Name,
		a.Address,
		a.Province,
		a.City,
		a.ZipCode,
	}
}
