package initializers

import "api-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
