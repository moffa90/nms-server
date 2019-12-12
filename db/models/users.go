package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"github.com/moffa90/nms-server/assets"
)

type User struct {
	Id       	string  `gorm:"not null;primary_key:true"`
	Username	string  `gorm:"not null;unique"`
	Name 		string  `gorm:"not null"`
	Password  	string	`gorm:"not null"`
	Role     	Role
	RoleID		string
	Active    	bool    `sql:"DEFAULT:true"`
	LastLogin	time.Time
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	DeletedAt 	*time.Time
}

type UserJson struct {
	User
	Role  string `validate:"required,ascii" json:"role"`
}

func MigrateUser(db *gorm.DB) {
	db.Debug().AutoMigrate(&User{})
	SeedUser(db)
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	if (!scope.DB().Where(&User{Username: user.Username}).First(&User{}).RecordNotFound()) {
		return errors.New("email duplicated")
	}

	id,_ := uuid.NewV4()
	scope.SetColumn("ID", id.String())

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	scope.SetColumn("Password", passwordHash)

	return nil
}


//Insert new user in the database
func CreateUser(db *gorm.DB, user *UserJson) (string, error) {
	tx := db.Begin()

	result := tx.Create(&user.User)
	if result.Error != nil {
		tx.Rollback()
		return "", result.Error
	}

	var r []Role
	db.Find(&r)
	if err := tx.Where(&Role{Name: user.Role}).First(&r).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return "", err
	}

	tx.Model(&user.User).Association("Role").Append(r)
	tx.Commit()

	return user.Id, nil
}

//Seeder for users, read file and insert into database
func SeedUser(db *gorm.DB) {
	file, _ := assets.Asset("goAssets/seeds/users-seed.json")

	var u []UserJson
	json.Unmarshal(file, &u)
	for _, element := range u {
		_, error := CreateUser(db, &element)

		if error != nil {
			log.Println(error.Error())
			continue
		}
	}
}

func GetUserByUsername(db *gorm.DB, username string) (User, error){
	var user User
	if err := db.Preload("Role").Where(&User{Username: username}).First(&user).Error; err != nil{
		return User{}, err
	}else{
		return user, nil
	}
}

func GetUserById(db *gorm.DB, id string) (User, error){
	var user User
	if err := db.Preload("Role").Where(&User{Id: id}).First(&user).Error; err != nil{
		return User{}, err
	}else{
		return user, nil
	}
}

func GetUsers(db *gorm.DB) ([]User) {
	var user []User

	if err := db.Preload("Role").Find(&user).Error; err != nil{
		return nil
	}else{
		return user
	}
}

func ChangeStatusUser(db *gorm.DB, userId string) error{
	var user User
	if err := db.Where(&User{Id: userId}).First(&user).Error; err != nil{
		return err
	}else{
		user.Active = !user.Active
		if err := db.Save(&user).Error; err != nil{
			return err
		}
		return nil
	}
}


