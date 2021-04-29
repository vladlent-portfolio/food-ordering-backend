package config

import (
	"path/filepath"
	"runtime"
)

var MaxUploadFileSize = 1 << 20 // 1 MiB

// StaticDir shows absolute path to project's "static" directory.
var StaticDir = filepath.Join(PathToMain(), "/static")

// CategoriesImgDir shows absolute path for categories images.
var CategoriesImgDir = filepath.Join(StaticDir, "categories", "img")

// PathToMain returns absolute path to project root.
func PathToMain() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../")
}
