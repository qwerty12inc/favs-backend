package models

type Category struct {
	Name          string   `firestore:"name" json:"name"`
	Labels        []string `firestore:"labels" json:"labels"`
	NeedsPurchase bool     `firestore:"needs_purchase" json:"needsPurchase"`
}
