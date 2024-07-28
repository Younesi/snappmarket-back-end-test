package main

import (
	"fmt"
	"log"
	"myapp/api/router"
	"myapp/config"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	validatorUil "myapp/util/validator"

	gormlogger "gorm.io/gorm/logger"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

//  @title          MYAPP API
//  @version        1.0
//  @description    This is a sample RESTful API with a CRUD

//  @contact.name   Mahdi Younesi

// @host       localhost:8080
// @basePath   /v1
func main() {
	c := config.New()
	v := validatorUil.New()

	var logLevel gormlogger.LogLevel
	if c.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatal("DB connection start failure")
		return
	}

	r := router.New(db, v)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server " + s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
