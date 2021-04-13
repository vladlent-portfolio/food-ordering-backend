package category

type DTO struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title"`
	Removable bool   `json:"removable"`
}
