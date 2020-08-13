package productsvc

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	errNotFound = errors.New("Product not found")
)

// Service represents methods exposed by the product service.
type Service interface {
	GetProduct(ctx context.Context, id string) (Product, error)
	PostProduct(ctx context.Context, product Product) (string, error)
	ListProduct(ctx context.Context) []Product
}

// Product is an item we sell.
type Product struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float32    `json:"price"`
	Categories  []Category `json:"categories"`
}

// Category represents a grouping of products.
type Category struct {
	Name        string `json:"name"`
	HumanName   string `json:"humanName"`
	Description string `json:"description"`
}

type inMemoryService struct {
	mtx        sync.RWMutex
	products   map[string]Product
	categories map[string]Category
}

// NewInMemoryService creates an implementation of service which is stored in
// memory rather than in a database.
func NewInMemoryService() Service {
	return &inMemoryService{
		products:   map[string]Product{},
		categories: map[string]Category{},
	}
}

func (svc *inMemoryService) GetProduct(ctx context.Context, id string) (Product, error) {
	svc.mtx.RLock()
	defer svc.mtx.Unlock()
	product, ok := svc.products[id]
	if !ok {
		return Product{}, errNotFound
	}
	return product, nil
}

func (svc *inMemoryService) PostProduct(ctx context.Context, product Product) (string, error) {
	svc.mtx.Lock()
	defer svc.mtx.Unlock()
	id := uuid.New().String()
	product.ID = id
	svc.products[id] = product
	return id, nil
}

func (svc *inMemoryService) ListProduct(ctx context.Context) []Product {
	results := make([]Product, 0)
	for _, product := range svc.products {
		results = append(results, product)
	}
	return results
}
