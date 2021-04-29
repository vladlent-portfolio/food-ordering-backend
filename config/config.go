package config

import (
	"os"
	"path/filepath"
)

var MaxUploadFileSize = 1 << 20 // 1 MiB

var CategoriesImgDir = filepath.Join(PathToStatic(), "categories", "img")

func PathToStatic() string {
	main, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Join(main, "/static")
}
