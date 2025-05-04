package dto

type Order struct {
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Description string  `json:"description"`
	Product     Product `json:"product"`
	Status      int     `json:"status"`
}

type Product struct {
	Name        string `json:"name"`
	Weigth      string `json:"weigth"`
	Description string `json:"description"`
}
