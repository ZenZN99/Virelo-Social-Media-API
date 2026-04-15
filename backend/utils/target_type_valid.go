package utils

import "backend/models"

func IsValidTargetType(t models.ContentType) bool {
	switch t {
	case models.Post, models.Reel, models.Comment:
		return true
	default:
		return false
	}
}
func IsValidLikeType(t models.ContentType) bool {
	switch t {
	case models.Post, models.Reel, models.Comment:
		return true
	default:
		return false
	}
}
