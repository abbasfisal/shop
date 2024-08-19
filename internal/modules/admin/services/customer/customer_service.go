package customer

import (
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/repositories/customer"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type CustomerService struct {
	repo customer.CustomerRepositoryInterface
}

func NewCustomerService(customerRepo customer.CustomerRepositoryInterface) CustomerService {
	return CustomerService{repo: customerRepo}
}

func (cs CustomerService) Index(c *gin.Context) (responses.Customers, custom_error.CustomError) {
	customers, err := cs.repo.GetAll(c)
	if err != nil {
		return responses.Customers{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}

	return responses.ToCustomers(customers), custom_error.CustomError{}
}
