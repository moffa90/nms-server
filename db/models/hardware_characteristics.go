package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type HardwareCharacteristics struct {
	HardID		string		`gorm:"not null;primary_key:true"`
	Serial		string		`gorm:"not null"`
	Key 		string		`gorm:"not null;primary_key:true"`
	Value 		string		`gorm:"not null"`
	Category	string		`gorm:"not null;primary_key:true"`
	Device     	Hardware
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewHardwareCharacteristics(hardID string, serial string, key string, value string, category string, device Hardware) *HardwareCharacteristics {
	return &HardwareCharacteristics{HardID: hardID, Serial: serial, Key: key, Value: value, Category: category, Device: device}
}

func MigrateHardwareCharacteristics(db *gorm.DB) {
	db.Debug().AutoMigrate(&HardwareCharacteristics{})
}

func (hChar *HardwareCharacteristics) BeforeCreate(scope *gorm.Scope) error {
	if (!scope.DB().Where(&HardwareCharacteristics{Serial: hChar.Serial, Key: hChar.Key, Category: hChar.Category}).First(&HardwareCharacteristics{}).RecordNotFound()) {
		return errors.New("Combination dev-key duplicated")
	}

	return nil
}

func (hChar *HardwareCharacteristics)Save(db *gorm.DB) error {

	var auxChar HardwareCharacteristics
	//log.Printf("Cat: %s, Key:%s", hChar.Category, hChar.Key)
	if (db.Where(&HardwareCharacteristics{HardID: hChar.HardID, Key: hChar.Key, Category: hChar.Category}).First(&auxChar).RecordNotFound()) {
		log.Printf("Not found ")
		//tx := db.Begin()

		result := db.Create(hChar)
		if result.Error != nil {
			return result.Error
		}

		//tx.Model(*hChar).Association("Device").Append(hChar.Device)
		//tx.Commit()
	}else{
		auxChar.Value = hChar.Value
		return db.Save(&auxChar).Error
	}

	return nil
}

func GetHardwareCharacteristicsByCat(serial, category string, db *gorm.DB) ([]HardwareCharacteristics, error){
	var chars []HardwareCharacteristics
	if err := db.Set("gorm:auto_preload", true).Where(&HardwareCharacteristics{Category: category, Serial: serial}).Find(&chars).Error; err != nil{
		return nil, err
	}else{
		return chars, nil
	}
}