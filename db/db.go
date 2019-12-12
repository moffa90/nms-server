package db

import (
	"github.com/moffa90/nms-server/db/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"sync"
)

var Shared *gorm.DB
var once sync.Once

func Init(){
	var err error
	//TODO: remove strings and replace with constants
	if _, err := os.Stat("/data/nms-server/nms-data"); os.IsNotExist(err) {
		os.OpenFile("/data/nms-server/nms-data", os.O_RDWR|os.O_CREATE, 0666)
	}


	once.Do(func() {
		Shared, err = gorm.Open("sqlite3", "file:/data/nms-server/nms-data?cache=shared&mode=rwc&_sync=1&_journal=WAL")
		//Shared, err = gorm.Open("sqlite3", "file:/data/nms-server/nms-data")
		if err != nil {
			panic("failed to connect database")
		}
	})


	models.MigrateRole(Shared)
	models.MigrateUser(Shared)
	models.MigrateInterface(Shared)
	models.MigrateHardware(Shared)
	models.MigrateConfiguration(Shared)
	//Shared.DropTable(models.HardwareCharacteristics{})
	models.MigrateHardwareCharacteristics(Shared)
	//Shared.DropTable(models.Remote{})
	models.MigrateRemotes(Shared)
}