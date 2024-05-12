package models

type Token struct {
	Token string `json:"token" firestore:"token"`
}
