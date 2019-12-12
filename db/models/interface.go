package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
	"github.com/moffa90/triadNMS/assets"
)

type Interface struct {
	Id        string     `gorm:"not null;primary_key:true"`
	Name      string     `gorm:"not null;unique"`
	DHCP	  bool		 `gorm:"not null;"`
	IPAddress string
	SubnetMask string
	Gateway   string
	SSID      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (i *Interface) BeforeCreate(scope *gorm.Scope) error {
	id,_ := uuid.NewV4()
	scope.SetColumn("id", id.String())
	return nil
}

func MigrateInterface(db *gorm.DB) {
	db.Debug().AutoMigrate(	&Interface{})
	SeedInterface(db)
}

func CreateInterface(db *gorm.DB, i *Interface) (string, error) {
	result := db.Create(i)
	if result.Error != nil {
		return "", result.Error
	}
	return i.Id, nil
}

func SeedInterface(db *gorm.DB) {
	file, _ := assets.Asset("goAssets/seeds/interfaces-seed.json")

	type interfaces struct {
		Name string
		Dhcp bool
	}

	var r []interfaces
	json.Unmarshal(file, &r)
	for _, element := range r {
		CreateInterface(db, &Interface{
			Name: element.Name,
			DHCP: element.Dhcp,
		})
	}
}

func GetInterfaceByName(db *gorm.DB, iName string) (Interface, error){
	var i Interface
	if err := db.Where(&Interface{Name: iName}).First(&i).Error; err != nil{
		return Interface{}, err
	}else{
		return i, nil
	}
}
