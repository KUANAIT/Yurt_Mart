package usecase

import (
	"log"
	natsPublisher "product-service/internal/nats"
	"time"

	"product-service/internal/cache"
	"product-service/internal/domain"
	"product-service/internal/repository"
)

type ProductUsecase interface {
	Create(product *domain.Product) (string, error)
	GetProductsByCategory(category string) ([]*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	Delete(id string) error
	List() ([]*domain.Product, error)
}

type productUsecase struct {
	repo      repository.ProductRepository
	cache     *cache.ProductCache
	publisher *natsPublisher.Publisher
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	c := cache.NewProductCache()

	products, err := repo.List()
	if err == nil {
		log.Println("✅ Product cache initialized on startup with", len(products), "items.")
		c.SetProducts(products)
	}

	p := natsPublisher.NewPublisher("nats://localhost:4222")

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			if c.ShouldRefresh() {
				log.Println("🔄 Refreshing product cache...")
				products, err := repo.List()
				if err == nil {
					c.SetProducts(products)
					log.Println("✅ Product cache refreshed.")
				}
			}
		}
	}()

	return &productUsecase{
		repo:      repo,
		cache:     c,
		publisher: p,
	}
}

func (u *productUsecase) Create(product *domain.Product) (string, error) {
	id, err := u.repo.Create(product)
	if err != nil {
		log.Println("❌ Failed to create product:", err)
		return "", err
	}
	product.ID = id

	u.cache.AddProduct(product)
	log.Println("🆕 Product added to cache with ID:", id)

	u.publisher.PublishProductCreated(product)

	return id, nil
}

func (u *productUsecase) Delete(id string) error {
	log.Println("🗑️ Deleting product with ID:", id)
	err := u.repo.Delete(id)
	if err != nil {
		return err
	}

	u.publisher.PublishProductDeleted(id)

	return nil
}

func (u *productUsecase) GetProductsByCategory(category string) ([]*domain.Product, error) {
	return u.repo.GetProductsByCategory(category)
}

func (u *productUsecase) List() ([]*domain.Product, error) {
	log.Println("📦 Returning products from cache...")
	return u.cache.GetAll(), nil
}

func (u *productUsecase) GetByID(id string) (*domain.Product, error) {
	if p, ok := u.cache.GetByID(id); ok {
		log.Println("📦 Returning product from cache with ID:", id)
		return p, nil
	}
	log.Println("🔎 Product not found in cache, fetching from DB:", id)
	return u.repo.GetByID(id)
}
