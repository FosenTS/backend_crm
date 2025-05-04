package model

type OrderStatus int8

const (
	Consideration = iota
	Refected
	AtWork
	Complete
)

type Order struct {
	OrderId     string
	Phone       string
	Email       string
	Description string
	Product     Product
	Status      OrderStatus
}
