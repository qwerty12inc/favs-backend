package models

type UserPurchases struct {
	UserEmail string           `json:"userEmail"`
	Objects   []PurchaseObject `json:"objects"`
}

func (up *UserPurchases) HasPurchase(id string) bool {
	for _, obj := range up.Objects {
		if obj.ID == id {
			return true
		}
	}
	return false
}

type PurchaseObject struct {
	ID    string `json:"id"` // stripe product id
	Price int    `json:"price"`
}
