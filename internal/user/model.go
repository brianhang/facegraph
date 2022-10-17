package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name string

	GoogleID string `gorm:"uniqueIndex:idx_user_google_id"`
}
