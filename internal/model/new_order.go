package model

type NewOrder struct {
	Phone       string
	Email       string
	Description string
	ProductId   string
	Status      OrderStatus
}
