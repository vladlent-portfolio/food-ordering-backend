package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	if IsProdMode {
		ex, err := os.Executable()

		if err != nil {
			panic(err)
		}

		exDir := filepath.Dir(ex)
		viper.SetConfigFile(filepath.Join(exDir, ".production.env"))
	} else {
		viper.SetConfigFile(filepath.Join(PathToMain(), ".env"))
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	HostRaw = viper.GetString("HOST_URL")
	parsedURL, err := url.Parse(HostRaw)

	if err != nil {
		fmt.Println("error parsing host url:", err)
	}

	HostURL = parsedURL
	ClientURL, _ = url.Parse(viper.GetString("FE_URL"))
}

var IsProdMode = gin.Mode() == "release"

var HostRaw string
var HostURL *url.URL
var ClientURL *url.URL

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
