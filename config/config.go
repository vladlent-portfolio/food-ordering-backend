package config

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/url"
	"path/filepath"
	"runtime"
)

func init() {
	envFileName := ".env"

	if IsProdMode {
		envFileName = ".env.production"
	}

	viper.SetConfigFile(envFileName)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

var IsProdMode = gin.Mode() == "release"

var HostRaw = "http://localhost:8080"
var HostURL, _ = url.Parse(HostRaw)

var MaxUploadFileSize int64 = 512 * 1024 // 512 KiB

// StaticDir shows path to "static" directory relative to main.go
var StaticDir = "static"

// CategoriesImgDir shows path to directory that contains categories
// images relative to static folder
var CategoriesImgDir = filepath.Join(StaticDir, "categories/img")

// DishesImgDir shows path to directory that contains dishes images
// relative to static folder
var DishesImgDir = filepath.Join(StaticDir, "dishes/img")

// StaticDirAbs shows absolute path to project's "static" directory.
var StaticDirAbs = filepath.Join(PathToMain(), StaticDir)

// CategoriesImgDirAbs shows absolute path for categories images.
var CategoriesImgDirAbs = filepath.Join(PathToMain(), CategoriesImgDir)

// DishesImgDirAbs shows absolute path for categories images.
var DishesImgDirAbs = filepath.Join(PathToMain(), DishesImgDir)

// PathToMain returns absolute path to project root.
func PathToMain() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../")
}
