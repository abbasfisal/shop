package customer

import (
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/repositories"

	"github.com/gin-gonic/gin"
)

type CustomerService struct {
	repo repositories.CustomerRepositoryInterface
}

func NewCustomerService(customerRepo repositories.CustomerRepositoryInterface) CustomerServiceInterface {
	return &CustomerService{repo: customerRepo}
}

func (cs *CustomerService) Index(c *gin.Context) (*responses.Customers, domain_err.CustomError) {
	customers, err := cs.repo.GetAll(c)
	if err != nil || customers == nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return responses.ToCustomers(customers), domain_err.CustomError{}
}
