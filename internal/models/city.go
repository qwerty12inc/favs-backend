package models

type City struct {
	Name       string      `firestore:"name" json:"name"`
	Center     Coordinates `firestore:"center" json:"center"`
	ImageURL   string      `firestore:"image_url" json:"imageURL"`
	Categories []Category  `firestore:"categories" json:"categories"`
}

type Category struct {
	Name    string   `firestore:"name" json:"name"`
	Filters []string `firestore:"filters" json:"filters"`
}
