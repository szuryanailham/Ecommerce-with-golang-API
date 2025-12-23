package products

type createProductRequest struct {
	Name       string `json:"name"`
	PriceCents int32  `json:"priceCents"`
	Quantity   int32  `json:"quantity"`
}

type UpdateProductNameRequest struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}