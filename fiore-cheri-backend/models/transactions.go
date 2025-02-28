package models

type Transaction struct {
	ID         int     `json:"id"`
	TotalPrice float64 `json:"total_price"`
	Timestamp  string  `json:"timestamp"`
}

type TransactionDetail struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	MenuID        int     `json:"menu_id"`
	Qty           int     `json:"qty"`
	Subtotal      float64 `json:"subtotal"`
}
