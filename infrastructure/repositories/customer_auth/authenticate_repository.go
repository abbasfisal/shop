package customer_auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shop/domain/entities"
)

type AuthenticateRepository struct {
	db *gorm.DB
}

func NewAuthenticateRepository(db *gorm.DB) AuthenticateRepository {
	return AuthenticateRepository{
		db: db,
	}
}

func (ar AuthenticateRepository) FindCustomerBySessionID(c *gin.Context, sessionID string) (entities.Customer, error) {
	var sess entities.Session

	err := ar.db.
		WithContext(c).
		Preload("Customer.Address").
		Preload("Customer.Carts.CartItems").
		Where("session_id = ?", sessionID).
		First(&sess).
		Error

	if err != nil {
		return entities.Customer{}, err
	}

	ar.loadCartItemAttributes(c, &sess.Customer)

	return sess.Customer, nil
}

// loadCartItemAttributes fills every cart line with the labels of the variant
// it points at, so /checkout/cart shows the exact combination the customer
// picked («رنگ: قرمز», «سایز: M») and not just the product title. One batched
// query for the whole cart — no extra round trip per line.
func (ar AuthenticateRepository) loadCartItemAttributes(c *gin.Context, customer *entities.Customer) {
	type variantRow struct {
		VariantID uint
		Title     string
		Value     string
	}

	variantIDs := []uint{}
	for _, cart := range customer.Carts {
		for _, item := range cart.CartItems {
			if item.InventoryID > 0 {
				variantIDs = append(variantIDs, item.InventoryID)
			}
		}
	}
	if len(variantIDs) == 0 {
		return
	}

	var rows []variantRow
	err := ar.db.
		WithContext(c).
		Table("variant_attribute_values AS vav").
		Select("vav.variant_id AS variant_id, av.attribute_title AS title, av.value AS value").
		Joins("JOIN attribute_values av ON av.id = vav.attribute_value_id AND av.deleted_at IS NULL").
		Where("vav.deleted_at IS NULL AND vav.variant_id IN ?", variantIDs).
		Order("av.attribute_id, av.sort_order").
		Scan(&rows).Error
	if err != nil {
		fmt.Println("[customer_auth]-[loadCartItemAttributes]-err:", err)
		return
	}

	byVariant := make(map[uint][]entities.CartItemAttribute, len(variantIDs))
	for _, row := range rows {
		byVariant[row.VariantID] = append(byVariant[row.VariantID], entities.CartItemAttribute{
			Title: row.Title,
			Value: row.Value,
		})
	}

	for ci := range customer.Carts {
		cart := &customer.Carts[ci]
		for ii := range cart.CartItems {
			item := &cart.CartItems[ii]
			item.Attributes = byVariant[item.InventoryID]
		}
	}
}
