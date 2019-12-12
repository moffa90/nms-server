package models

import (
	"github.com/moffa90/nms-server/assets"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type Role struct {
	Id        string     `gorm:"not null;primary_key:true"`
	Name      string     `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (r *Role) BeforeCreate(scope *gorm.Scope) error {
	id,_ := uuid.NewV4()
	scope.SetColumn("id", id.String())
	return nil
}

func MigrateRole(db *gorm.DB) {
	db.Debug().AutoMigrate(&Role{})
	SeedRole(db)
}

func CreateRole(db *gorm.DB, role *Role) (string, error) {
	result := db.Create(role)
	if result.Error != nil {
		return "", result.Error
	}
	return role.Id, nil
}

func SeedRole(db *gorm.DB) {
	file, _ := assets.Asset("goAssets/seeds/roles-seed.json")

	type roles struct {
		Name string
	}

	var r []roles
	json.Unmarshal(file, &r)
	for _, element := range r {
		CreateRole(db, &Role{
			Name: strings.ToLower(element.Name),
		})
	}
}

func GetRoles(db *gorm.DB) ([]Role){
	var r []Role

	if err := db.Find(&r).Error; err != nil{
		return nil
	}else{
		return r
	}
}