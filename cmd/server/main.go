package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/amaury95/GetGround-Party/api"
	"github.com/amaury95/GetGround-Party/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// setup parser
	parser := argparse.NewParser("party", "Party is the webserver to manage GetGround invitations and guests.")

	// setup parser arguments
	var (
		port     = parser.Int("p", "port", &argparse.Options{Default: 3033, Help: `server port to listen for requests`})
		username = parser.String("n", "username", &argparse.Options{Default: "root", Help: `mysql connection username`})
		password = parser.String("k", "password", &argparse.Options{Default: "example", Help: `mysql connection password`})
		url      = parser.String("u", "url", &argparse.Options{Default: "127.0.0.1:3306", Help: `mysql connection server url`})
		database = parser.String("d", "database", &argparse.Options{Default: "party", Help: `mysql connection database name`})
	)

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	// open db connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", *username, *password, *url, *database)
	db, err := gorm.Open(mysql.Open(dsn), new(gorm.Config))
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

	if err := router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		panic("error running the server: " + err.Error())
	}
}
