package category

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;unique;not null"`
	Removable bool
	Image     *string
}
