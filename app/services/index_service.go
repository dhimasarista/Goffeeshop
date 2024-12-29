package services

import (
	"Goffeeshop/app/models"
	"Goffeeshop/app/repositories"
	"log"
)

type IndexService struct {
	ProductRepo *repositories.ProductRepository
}

func NewIndexService(productRepo *repositories.ProductRepository) *IndexService {
	return &IndexService{ProductRepo: productRepo}
}

func (is *IndexService) Index() ([]models.Product, error) {
	products, err := is.ProductRepo.All()
	if err != nil {
		log.Println("Error fetching products:", err)
		return nil, err
	}
	return products, nil

}
