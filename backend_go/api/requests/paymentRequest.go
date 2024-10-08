package requests

type PaymentRequest struct {
	Type     string  `json:"type" binding:"required"`
	Quantity uint    `json:"quantity" binding:"required"`
	Price    float32 `json:"price" binding:"required"`
}
