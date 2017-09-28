package md

import "github.com/jinzhu/gorm"


type Page struct {
	gorm.Model
	Name string
	Path string
	FirstLine string
	Tags []string
}
