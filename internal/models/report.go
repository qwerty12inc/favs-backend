package models

type Report struct {
	ID          string `json:"id,omitempty" firestore:"id"`
	PlaceID     string `json:"place_id,omitempty" firestore:"place_id"`
	ReportedBy  string `json:"reported_by" firestore:"reported_by"`
	ReportedAt  int64  `json:"reported_at,omitempty" firestore:"reported_at"`
	Description string `json:"description" firestore:"description"`
}
