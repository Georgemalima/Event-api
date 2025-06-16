package store

type EventCustomer struct {
	ID         int64  `json:"id"`
	EventID    int64  `json:"event_id"`
	CustomerID int64  `json:"customer_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
