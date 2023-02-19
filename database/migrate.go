package database

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/models"
	"fmt"
)

func CreateAutoMigration(database *controller.Database) {
	fmt.Println("running migrations................................")
	database.DB.AutoMigrate(&models.User{})
	//database.DB.Migrator().DropTable(&User{})

}
