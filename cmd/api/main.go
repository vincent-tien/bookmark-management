package main

import (
	_ "github.com/vincent-tien/bookmark-management/docs"
	"github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/pkg/logger"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
)

// @title	Bookmark Management API
// @version	1.0
// @description	Bookmark Management API is a RESTful service
func main() {
	logger.SetLogLevel()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	app := api.New(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
