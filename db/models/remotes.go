package models

import (
	uuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Remote struct {
	Id 			string  `gorm:"primary_key;auto_increment:false"`
	Remote		int8 	`gorm:"unique_index:idx_name_code"`
	Group		int8 	`gorm:"unique_index:idx_name_code"`
	Ip			string
	Port		string
	Hostname    string
	Alarms		string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (r *Remote) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.New()
	scope.SetColumn("id", id.String())
	return nil
}

func CreateRemote(db *gorm.DB, r *Remote) error {
	result := db.Create(r)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func MigrateRemotes(db *gorm.DB) {
	db.Debug().AutoMigrate(&Remote{})
}

func CreateOrUpdate(db *gorm.DB, r *Remote) error{
	var aux Remote

	if db.Where("remote == ? AND `group` == ?", r.Remote, r.Group).First(&aux).RecordNotFound() {
		return CreateRemote(db, r)
	}else{
		aux.Ip = r.Ip
		aux.Port = r.Port
		if r.Hostname != "" {
			aux.Hostname = r.Hostname
		}
		return db.Save(&aux).Error
	}
}

func GetRemotes(db *gorm.DB) []Remote{
	var r []Remote

	if err := db.Find(&r).Error; err != nil{
		return nil
	}else{
		return r
	}
}