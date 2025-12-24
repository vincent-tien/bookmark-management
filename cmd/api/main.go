package main

import (
	_ "github.com/vincent-tien/bookmark-management/docs"
	"github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
)

// @title	Bookmark Management API
// @version	1.0
// @description	Bookmark Management API is a RESTful service
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	app := api.New(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
