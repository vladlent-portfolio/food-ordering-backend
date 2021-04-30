package category

import (
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"path"
)

// PathToImg returns absolute URI from config.HostURL to provided image
// in config.CategoriesImgDir folder.
func PathToImg(imgName string) string {
	return common.HostURLResolver(
		path.Join(config.CategoriesImgDir, imgName),
	)
}
