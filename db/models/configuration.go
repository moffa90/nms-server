package models

import (
	"github.com/moffa90/triadNMS/assets"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type Configuration struct {
	Id      string  `gorm:"not null;primary_key:true"`
	Key		string  `gorm:"not null;unique"`
	Value	string	`gorm:"not null;"`
	Category string `gorm:"not null;default:general"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func MigrateConfiguration(db *gorm.DB) {
	db.Debug().AutoMigrate(	&Configuration{})
	SeedConfiguration(db)
}

func (c *Configuration) BeforeCreate(scope *gorm.Scope) error {
	id,_ := uuid.NewV4()
	scope.SetColumn("id", id.String())
	return nil
}

func CreateConfiguration(db *gorm.DB, c *Configuration) (string, error) {
	result := db.Create(c)
	if result.Error != nil {
		return "", result.Error
	}
	return c.Id, nil
}

//Seeder for users, read file and insert into database
func SeedConfiguration(db *gorm.DB) {
	file, _ := assets.Asset("goAssets/seeds/configuration-seed.json")

	var c []Configuration
	json.Unmarshal(file, &c)
	for _, element := range c {
		_, error := CreateConfiguration(db, &element)

		if error != nil {
			log.Println(error.Error())
			continue
		}
	}
}

func GetConfigByKey(k string, db *gorm.DB) (*Configuration, error){
	if k == ""{
		return  nil, errors.New("Empty Key")
	}

	var conf Configuration
	if err := db.Where(&Configuration{Key:k}).First(&conf).Error; err != nil {
		return nil,err
	}

	return &conf, nil
}

func GetConfigByCat(c string, db *gorm.DB) ([]Configuration, error){
	if c == ""{
		return  nil, errors.New("empty Key")
	}

	var conf []Configuration
	if err := db.Where(&Configuration{Category:c}).Find(&conf).Error; err != nil {
		return nil,err
	}

	return conf, nil
}

func GetConfig(db *gorm.DB) ([]Configuration, error){
	var conf []Configuration
	if err := db.Find(&conf).Error; err != nil {
		return nil, err
	}

	return conf, nil
}