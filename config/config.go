package config

import (
	"net/url"
	"path/filepath"
	"runtime"
)

const HostRaw = "http://localhost:8080"

var HostURL, _ = url.Parse(HostRaw)

var MaxUploadFileSize int64 = 512 * 1024 // 512 KiB

// StaticDir shows path to "static" directory relative to main.go
var StaticDir = "static"

// CategoriesImgDir shows path to categories images directory relative to main.go
var CategoriesImgDir = "static/categories/img"

// StaticDirAbs shows absolute path to project's "static" directory.
var StaticDirAbs = filepath.Join(PathToMain(), StaticDir)

// CategoriesImgDirAbs shows absolute path for categories images.
var CategoriesImgDirAbs = filepath.Join(PathToMain(), CategoriesImgDir)

// PathToMain returns absolute path to project root.
func PathToMain() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../")
}
