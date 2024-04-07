package stripe

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/product"
)

type StripeConnector interface {
	GetProducts() ([]Product, error)
	GetProductByID(id string) (Product, error)
	GetPricesByProductID(id string) (uint, error)
}

type Product struct {
	ProductID string
	PriceID   string
	Price     int64
}

type StripeConnectorImpl struct {
}

func NewStripeConnector() *StripeConnectorImpl {
	return &StripeConnectorImpl{}
}

func (s *StripeConnectorImpl) GetProducts() ([]Product, error) {
	products := make([]Product, 0)

	productParams := &stripe.ProductListParams{}
	productIterator := product.List(productParams)
	for productIterator.Next() {
		products = append(products, Product{
			ProductID: productIterator.Product().ID,
			PriceID:   productIterator.Product().DefaultPrice.ID,
			Price:     productIterator.Product().DefaultPrice.UnitAmount,
		})
	}

	return products, nil
}

func (s *StripeConnectorImpl) GetProductByID(id string) (Product, error) {
	productParams := &stripe.ProductParams{
		ID: stripe.String(id),
	}
	pr, err := product.Get(id, productParams)
	if err != nil {
		return Product{}, err
	}

	return Product{
		ProductID: pr.ID,
		PriceID:   pr.DefaultPrice.ID,
		Price:     pr.DefaultPrice.UnitAmount,
	}, nil
}

func (s *StripeConnectorImpl) GetPricesByProductID(id string) (uint, error) {
	productParams := &stripe.ProductParams{
		ID: stripe.String(id),
	}
	pr, err := product.Get(id, productParams)
	if err != nil {
		return 0, err
	}

	return uint(pr.DefaultPrice.UnitAmount), nil
}
