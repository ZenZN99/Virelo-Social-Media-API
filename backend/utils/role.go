package utils

import "backend/models"

func GetRole() models.UserRole {

	return models.UserRole("user")
}
