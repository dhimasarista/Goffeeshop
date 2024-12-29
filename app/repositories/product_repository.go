package repositories

import (
	"Goffeeshop/app/models"
	"log"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) All() ([]models.Product, error) {
	var products []models.Product

	err := repo.DB.Find(&products).Error
	if err != nil {
		log.Println("Error get all products: ", err)
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) First(id string) (models.Product, error) {
	var product models.Product
	err := repo.DB.First(&product, "id = ?", id).Error
	if err != nil {
		log.Println("Error find product: ", err)
		return models.Product{}, err
	}

	return product, nil
}
