package main

import (
	"github.com/amaury95/GetGround-Party/api"
	"github.com/amaury95/GetGround-Party/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// open db connection
	db, err := gorm.Open(mysql.Open("root:example@tcp(127.0.0.1:3306)/party?charset=utf8mb4&parseTime=true"), new(gorm.Config))
	if err != nil {
		panic("error connecting to database: " + err.Error())
	}

	// migrate the models to create database tables
	db.AutoMigrate(new(models.Table))
	db.AutoMigrate(new(models.Guest))
	db.AutoMigrate(new(models.Reservation))

	// create router instance
	router := new(api.Handler).WithConnection(db).Router(&api.RouterConfig{
		ShowLogs:    true,
		ReleaseMode: true,
	})

	if err := router.Run(`:3000`); err != nil {
		panic("error running the server: " + err.Error())
	}
}
