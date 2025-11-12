package dto

type CreateCatalogStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
