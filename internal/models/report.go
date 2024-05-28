package models

type Report struct {
	ID          string `json:"id",firestore:"id"`
	PlaceID     string `json:"place_id",firestore:"place_id"`
	ReportedBy  string `json:"reported_by",firestore:"reported_by"`
	ReportedAt  int64  `json:"reported_at",firestore:"reported_at"`
	Description string `json:"description",firestore:"description"`
}
