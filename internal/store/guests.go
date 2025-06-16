package store

type Guest struct {
	ID          int64  `json:"id"`
	Name        string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
	CardID      int64  `json:"card_id"`
	EventID     int64  `json:"role_id"`
	Event       Event  `json:"event"`
}
