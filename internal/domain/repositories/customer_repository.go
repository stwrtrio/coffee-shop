package repositories

import (
	"context"

	"github.com/stwrtrio/coffee-shop/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) error
	FindCustomerByEmail(ctx context.Context, email string) (*models.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateCustomer(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) FindCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &customer, err
}
