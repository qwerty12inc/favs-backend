package models

type Place struct {
	ID             string          `firestore:"id" json:"id"`
	Name           string          `firestore:"name" json:"name"`
	Category       string          `firestore:"category" json:"category"`
	Description    string          `firestore:"description" json:"description"`
	LocationURL    string          `firestore:"location_url" json:"locationURL"`
	Coordinates    Coordinates     `firestore:"coordinates" json:"coordinates"`
	City           string          `firestore:"city" json:"city"`
	Website        string          `firestore:"website" json:"website"`
	Instagram      string          `firestore:"instagram" json:"instagram"`
	Labels         []string        `firestore:"labels" json:"labels"`
	GeoHash        string          `firestore:"geohash" json:"geoHash"`
	Address        string          `firestore:"address" json:"address"`
	GoogleMapsInfo *GoogleMapsInfo `firestore:"google_maps_info" json:"googleMapsInfo"`
}

type GoogleMapsInfo struct {
	PlaceID          string   `firestore:"place_id" json:"placeID"`
	OpeningInfo      []string `firestore:"opening_info" json:"openingInfo"`
	FormattedAddress string   `firestore:"formatted_address" json:"formattedAddress"`
	LocationURL      string   `firestore:"location_url" json:"locationURL"`
	Rating           float32  `firestore:"rating" json:"rating"`
	Website          string   `firestore:"website" json:"website"`
	Reservable       bool     `firestore:"reservable" json:"reservable"`
	Delivery         bool     `firestore:"delivery" json:"delivery"`
	PhotoRefList     []string `firestore:"photo_ref_list" json:"photoRef"`
}

type Coordinates struct {
	Latitude  float64 `firestore:"latitude" json:"latitude"`
	Longitude float64 `firestore:"longitude" json:"longitude"`
}

type GoogleSheetPlace struct {
	Name        string
	Description string
	LocationURL string
	City        string
	Website     string
	Instagram   string
	Labels      []string
	Category    string
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
	City           string      `json:"city,omitempty"`
	Center         Coordinates `json:"center,omitempty"`
	LatitudeDelta  float64     `json:"latitude_delta,omitempty"`
	LongitudeDelta float64     `json:"longitude_delta,omitempty"`
	Labels         []string    `json:"labels,omitempty"`
}

type GetFiltersRequest struct {
	City string `json:"city,omitempty"`
}
