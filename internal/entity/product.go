package entity

import (
	"errors"
	"time"

	"github.com/MatheusFontesDaGama/api-go/pkg/entity"
)

var (
	ErrorIDIsRequired    = errors.New("id is required")
	ErrorInvalidID       = errors.New("invalid id")
	ErrorNameIsRequired  = errors.New("name is required")
	ErrorPriceIsRequired = errors.New("price is required")
	ErrorInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func (product *Product) Validate() error {
	if product.ID.String() == "" {
		return ErrorIDIsRequired
	}

	if _, err := entity.ParseID(product.ID.String()); err != nil {
		return ErrorInvalidID
	}

	if product.Name == "" {
		return ErrorNameIsRequired
	}

	if product.Price == 0 {
		return ErrorPriceIsRequired
	}

	if product.Price < 0 {
		return ErrorInvalidPrice
	}

	return nil
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil
}
