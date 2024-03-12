package models

type Place struct {
	ID          string      `firestore:"id"`
	Name        string      `firestore:"name"`
	Description string      `firestore:"description"`
	LocationURL string      `firestore:"location_url"`
	Coordinates Coordinates `firestore:"location"`
	City        string      `firestore:"city"`
	Website     string      `firestore:"website"`
	Instagram   string      `firestore:"instagram"`
	Labels      []string    `firestore:"labels"`
}

type Coordinates struct {
	Latitude  float64 `firestore:"lat"`
	Longitude float64 `firestore:"lng"`
}

type GoogleSheetPlace struct {
	Name        string
	Description string
	LocationURL string
	City        string
	Website     string
	Instagram   string
	Labels      []string
}

type CreatePlaceRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	LocationURL string   `json:"location,omitempty"`
	OpenAt      string   `json:"open_at,omitempty"`
	ClosedAt    string   `json:"closed_at,omitempty"`
	City        string   `json:"city,omitempty"`
	Website     string   `json:"website,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type UpdatePlaceRequest struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	LocationURL string   `json:"location,omitempty"`
	City        string   `json:"city,omitempty"`
	Address     string   `json:"address,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Type        string   `json:"type,omitempty"`
	Website     string   `json:"website,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

type GetPlacesRequest struct {
	Center         Coordinates `json:"center,omitempty"`
	LatitudeDelta  float64     `json:"latitude_delta,omitempty"`
	LongitudeDelta float64     `json:"longitude_delta,omitempty"`
	Labels         []string    `json:"labels,omitempty"`
}
