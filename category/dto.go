package category

type DTO struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Removable bool   `json:"removable"`
}
