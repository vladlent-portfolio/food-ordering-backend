package dish

import (
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"path"
	"path/filepath"
)

// PathToImg returns absolute URI from config.HostURL to provided image
// in config.DishesImgDir folder.
func PathToImg(imgName string) string {
	return common.HostURLResolver(
		// To make sure that the path is correct even if running on Microsoft Windows
		filepath.ToSlash(path.Join(config.DishesImgDir, imgName)),
	)
}
