package main

import (
	"github.com/pressly/goose/v3"
	_ "github.com/vincent-tien/bookmark-management/docs"
	"github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/pkg/logger"
	"github.com/vincent-tien/bookmark-management/pkg/migrations"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
	sqldbPkg "github.com/vincent-tien/bookmark-management/pkg/sqldb"
)

// @title	Bookmark Management API
// @version	1.0.2
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

	db, err := sqldbPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDB, "."); err != nil {
		panic(err)
	}

	app := api.New(cfg, redisClient, db)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
