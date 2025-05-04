package dto

type NewOrder struct {
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Description string `json:"description"`
	ProductId   string `json:"productId"`
}
