package models

type User struct {
	UID   string `firestore:"id"`
	Email string `firestore:"email"`
}
