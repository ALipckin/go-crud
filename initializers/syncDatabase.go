package initializers

import (
	"go-crud/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
		return
	}
}
