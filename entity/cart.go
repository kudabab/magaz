package entity

// товар в корзине
type CartItem struct {
	ProductID int     `json:"product_id"`
	Category  string  `json:"category"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

type Order struct {
	OrderID int        `json:"order_id"`
	UserID  int        `json:"user_id"`
	Items   []CartItem `json:"items"`
	Total   float64    `json:"total"`
	Status  string     `json:"status"`
}
