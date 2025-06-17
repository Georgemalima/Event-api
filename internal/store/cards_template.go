package store

type CardTemplate struct {
	ID        int64  `json:"id"`
	ImagePath string `json:"image_path"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
