package mygorm

import (
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/models"
	"github.com/varakornpz/providers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



var DB *gorm.DB

func InitDB(){
	db, err := gorm.Open(postgres.Open(providers.AppConf.DBDsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg("Cant connect to db")
	}

	DB = db

	db.AutoMigrate(&models.User{})
}
