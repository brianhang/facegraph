package user

import "brianhang.me/facegraph/internal/db"

func FindOrCreateFromGoogleID(googleID string) User {
	db := db.Get()

	var user User
	db.Where(User{GoogleID: googleID}).FirstOrCreate(&user, User{
		GoogleID: googleID,
	})
	return user
}
