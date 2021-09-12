package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type Hardware struct {
	Id        	string     	`gorm:"not null;primary_key:true"`
	DevId     	string     	`gorm:"not null;unique"`
	ProductId  	string		`gorm:"not null;"`
	Address  	int			`gorm:"not null;"`
	Serial		string		`gorm:"not null;DEFAULT:''"`
	Backplane	int			`gorm:"not null;DEFAULT:0"`
	Alarms		int
	FWversion	string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (h *Hardware) ToString() string{
	return fmt.Sprintf("ID=%s, DevId=%s, ProductId=%s, Slot Address=0x%02X, Serial=%s, Backplane=0x%02X, Alarms=0x%02X, Firmware Version=%s",
		h.Id, h.DevId, h.ProductId, h.Address, h.Serial, h.Backplane, h.Alarms, h.FWversion)
}

func (h *Hardware) Update(data map[string]string, db *gorm.DB) error{
	num, _ := strconv.ParseInt(strings.Replace(data["Alarms"], "x", "", -1), 16, 32)
	h.Alarms = int(num)
	num, _ = strconv.ParseInt(data["Slot"], 10, 32)
	h.Address = int(num)
	num, _ = strconv.ParseInt(data["Backplane"], 10, 32)
	h.Backplane = int(num)
	h.FWversion = data["Version"]
	return  db.Save(h).Error
}

func MigrateHardware(db *gorm.DB) {
	db.Debug().AutoMigrate(	&Hardware{})
}

func (h *Hardware) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.New()
	scope.SetColumn("id", id.String())
	return nil
}

func CreateHardware(db *gorm.DB, h *Hardware) (string, error) {
	result := db.Create(h)
	if result.Error != nil {
		return "", result.Error
	}
	return h.Id, nil
}

func GetHardware(db *gorm.DB) ([]Hardware, error){
	var h []Hardware
	if err := db.Find(&h).Error; err != nil{
		return nil, err
	}else{
		return h, nil
	}
}

func GetHardwareBySerial(db *gorm.DB, serial string) (*Hardware, error){
	var h Hardware
	if err := db.Where(Hardware{Serial:serial}).First(&h).Error; err != nil{
		return nil, err
	}else{
		return &h, nil
	}
}

func GetHardwareById(db *gorm.DB, id string) (*Hardware, error){
	var h Hardware
	if err := db.Where(Hardware{Id:id}).First(&h).Error; err != nil{
		return nil, err
	}else{
		return &h, nil
	}
}
