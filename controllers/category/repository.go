package category

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func ProvideRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(c Category) (Category, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *Repository) Save(c Category) (Category, error) {
	err := r.db.Save(&c).Error
	return c, err
}

func (r *Repository) FindByID(id uint) (Category, error) {
	var c Category
	err := r.db.First(&c, id).Error
	return c, err
}

func (r *Repository) FindAll() []Category {
	var categories []Category
	r.db.Order("id ASC").Find(&categories)
	return categories
}

func (r *Repository) FindAllDishImages(categoryID uint) ([]string, error) {
	var res []string
	tx := r.db.Table("dishes").Select("image").Where("image IS NOT NULL")

	if categoryID != 0 {
		tx.Where("category_id = ?", categoryID)
	}

	rows, err := tx.Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image string
		if err = rows.Scan(&image); err != nil {
			return nil, err
		}
		res = append(res, image)
	}

	return res, err
}

func (r *Repository) Delete(c Category) (Category, error) {
	err := r.db.Delete(&c).Error
	return c, err
}
