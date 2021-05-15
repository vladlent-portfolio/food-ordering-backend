package dish

import (
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"path/filepath"
)

// PathToImg returns absolute URI from config.HostURL to provided image
// in config.DishesImgDir folder.
func PathToImg(imgName string) string {
	return common.HostURLResolver(
		filepath.Join(config.DishesImgDir, imgName),
	)
}
