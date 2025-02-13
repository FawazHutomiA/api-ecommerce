package role

import (
	"github.com/google/uuid"
)

type RoleDetailResponse struct {
	ID          uuid.UUID `json:"ID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
