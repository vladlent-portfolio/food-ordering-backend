package dish

import (
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"path"
)

// PathToImg returns absolute URI from config.HostURL to provided image
// in config.DishesImgDir folder.
func PathToImg(imgName string) string {
	return common.HostURLResolver(
		path.Join(config.DishesImgDir, imgName),
	)
}
