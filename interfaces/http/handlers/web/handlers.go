package handlers

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	responses "shop/application/dto/admin"
	"shop/application/usecases/banner"
	"shop/application/usecases/fee"
	"shop/application/usecases/home"
	sliders "shop/application/usecases/product_slider"
	"shop/bootstrap"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/interfaces/http/response"
	"shop/pkg/helpers"
	"shop/pkg/util"
)

type PublicHandler struct {
	homeSrv    home.HomeServiceInterface
	bannerSrv  *banner.BannerService
	slidersSrv *sliders.ProductSliderService
	feeSrv     *fee.FeeRateService
	dep        *bootstrap.Dependencies
}

func NewPublicHandler(
	homeSrv home.HomeServiceInterface,
	bannerSrv *banner.BannerService,
	slidersSrv *sliders.ProductSliderService,
	feeSrv *fee.FeeRateService,
	dep *bootstrap.Dependencies,
) PublicHandler {
	return PublicHandler{
		homeSrv:    homeSrv,
		bannerSrv:  bannerSrv,
		slidersSrv: slidersSrv,
		feeSrv:     feeSrv,
		dep:        dep,
	}
}

func (p PublicHandler) Index(c *gin.Context) {
	//20 latest
	//20 random /-
	//20 lowest quantity
	//10 category
	//filter by price
	//filter by quantity
	//pagination in selected category

	products, err := p.homeSrv.GetProducts(context.TODO(), 20)
	if err.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if err.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	categories, cErr := p.homeSrv.GetCategories(context.TODO(), 20)
	if cErr.Code == 404 {
		c.JSON(200, gin.H{"err": 404})
		return
	}
	if cErr.Code == 500 {
		c.JSON(200, gin.H{"err": 500})
		return
	}

	c.JSON(200, gin.H{
		"PRODUCTS":   products,
		"CATEGORIES": categories,
	})
	//	response.Render(c, http.StatusFound, "home", gin.H{
	//			"TITLE":      "home page",
	//			"PRODUCTS":   "",
	//			"CATEGORIES": "",
	//		})
	return
}

func (p PublicHandler) ShowProductsByCategory(c *gin.Context) {

	productPagination, err := p.homeSrv.ListProductByCategorySlug(c, c.Param("category_slug"))
	if err != nil {
		//هر خطایی به جز خطای مرتبط با پیدانکردن رکورد اگر وجود داشت اون خطا رو نشون میدیم
		//در غیر این صورت پیغام رکورد یافت نشد به کاربر نشون داده میشه :)
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"msg": domain_err.SomethingWrongHappened,
			})
			return
		}
	}

	response.CustomerRender(c, http.StatusFound, "search",
		gin.H{
			"TITLE":          "search",
			"PAGINATION":     productPagination,
			"MEDIA_PATH":     util.GetProductStoragePath(),
			"PrimaryMessage": domain_err.RecordNotFound,
		},
	)
	return

}

// HomePage renders the storefront home: promotion banners (2-up / 4-up) and
// the curated product sliders, each pinned to its homepage slot.
func (p PublicHandler) HomePage(c *gin.Context) {
	bannersTwo := responses.Banners{}
	bannersFour := responses.Banners{}
	if active, err := p.bannerSrv.Active(c); err == nil {
		for _, b := range responses.ToBanners(active).Data {
			if b.Layout == entities.BannerLayoutFour {
				bannersFour.Data = append(bannersFour.Data, b)
			} else {
				bannersTwo.Data = append(bannersTwo.Data, b)
			}
		}
	}

	slidersByPosition := map[string][]responses.ProductSlider{}
	if active, err := p.slidersSrv.ActiveSliders(c); err == nil {
		for position, views := range responses.SlidersByPosition(active) {
			for i := range views {
				views[i].MediaPath = util.GetProductStoragePath()
			}
			slidersByPosition[position] = views
		}
	}

	response.CustomerRender(c, 200, "home", gin.H{
		"TITLE":        "صفحه اصلی فروشگاه",
		"BANNERS_TWO":  bannersTwo,
		"BANNERS_FOUR": bannersFour,
		"SLIDERS":      slidersByPosition,
		"BANNER_PATH":  util.GetBannerStoragePath(),
		"MEDIA_PATH":   util.GetProductStoragePath(),
	})
}

func (p PublicHandler) SingleProduct(c *gin.Context) {

	//get single product with its recommendations product
	product, recommendations, err := p.homeSrv.GetSingleProduct(c, c.Param("product_sku"), c.Param("product_slug"))
	primaryMessage := ""
	if err.Code > 0 {
		if err.Code == 404 {
			primaryMessage = err.DisplayMessage
		} else {
			response.Error500(c)
			return
		}
	}

	groups, variantsJSON := buildVariantPickerData(product)

	response.CustomerRender(c, http.StatusFound, "single_product",
		gin.H{
			"PRODUCT":         product,
			"RECOMMENDATIONS": recommendations,
			"PrimaryMessage":  primaryMessage,
			"MEDIA_PATH":      util.GetProductStoragePath(),
			"VARIANT_GROUPS":  groups,
			"VARIANTS_JSON":   variantsJSON,
			// variants of this product already sitting in the cart: the picker
			// swaps «افزودن به سبد» for «مشاهده سبد خرید» per combination
			"CART_INVENTORY_IDS_JSON": cartInventoryIDsJSON(c, product),
		})
	return
}

// cartInventoryIDsJSON serializes the inventory ids of this product that the
// logged-in customer already has in the cart (e.g. [2,7]).
func cartInventoryIDsJSON(c *gin.Context, product map[string]interface{}) string {
	ids := []uint{}
	if raw, ok := product["_id"].(string); ok {
		if productID, err := strconv.ParseUint(raw, 10, 64); err == nil {
			if customer, authed := helpers.GetAuthUser(c); authed {
				for _, item := range customer.Cart.CartItem.Data {
					if item.ProductID == uint(productID) {
						ids = append(ids, item.InventoryID)
					}
				}
			}
		}
	}
	payload, _ := json.Marshal(ids)
	return string(payload)
}

// VariantGroupValue is one selectable value inside a variant group.
type VariantGroupValue struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Hex   string `json:"hex"`
}

// VariantGroup is one attribute (e.g. رنگ) with its distinct values,
// rendered as swatches (color) or boxes (everything else).
type VariantGroup struct {
	AttributeID int64               `json:"attribute_id"`
	Title       string              `json:"title"`
	IsColor     bool                `json:"is_color"`
	Values      []VariantGroupValue `json:"values"`
}

// variantJSON is one sellable combination for the picker script.
type variantJSON struct {
	ID          int64   `json:"id"`
	ValueIDs    []int64 `json:"value_ids"`
	Effective   int64   `json:"effective"`
	Price       int64   `json:"price"`
	Percent     int64   `json:"percent"`
	HasDiscount bool    `json:"has_discount"`
	Available   int64   `json:"available"`
}

// buildVariantPickerData groups the read-model inventories by attribute
// (deterministic order) and serializes the combinations for the picker
// script. Stock-only rows (no attribute links) contribute combinations
// but no groups.
func buildVariantPickerData(product map[string]interface{}) ([]VariantGroup, string) {
	var inventories map[string]entities.Inventory
	if raw, ok := product["inventories"]; ok {
		inventories, _ = raw.(map[string]entities.Inventory)
	}

	// numeric inventory ids, ascending (map order is random)
	ids := make([]int64, 0, len(inventories))
	for key := range inventories {
		if id, err := strconv.ParseInt(key, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	groups := map[int64]*VariantGroup{}
	groupOrder := []int64{}
	variants := make([]variantJSON, 0, len(ids))

	for _, id := range ids {
		inv := inventories[strconv.FormatInt(id, 10)]
		valueIDs := make([]int64, 0, len(inv.Attributes))
		for _, attr := range inv.Attributes {
			if attr.AttributeID == 0 && attr.AttributeValueID == 0 {
				continue
			}
			valueIDs = append(valueIDs, attr.AttributeValueID)

			g, ok := groups[attr.AttributeID]
			if !ok {
				g = &VariantGroup{
					AttributeID: attr.AttributeID,
					Title:       attr.AttributeTitle,
					IsColor:     attr.IsColor,
				}
				groups[attr.AttributeID] = g
				groupOrder = append(groupOrder, attr.AttributeID)
			}
			if attr.IsColor {
				g.IsColor = true
			}
			seen := false
			for _, v := range g.Values {
				if v.ID == attr.AttributeValueID {
					seen = true
					break
				}
			}
			if seen {
				continue
			}
			hex := ""
			if attr.IsColor {
				hex = attr.ColorHex
			}
			g.Values = append(g.Values, VariantGroupValue{
				ID:    attr.AttributeValueID,
				Title: attr.AttributeValueTitle,
				Hex:   hex,
			})
		}

		variants = append(variants, variantJSON{
			ID:          inv.InventoryID,
			ValueIDs:    valueIDs,
			Effective:   inv.EffectivePrice,
			Price:       inv.Price,
			Percent:     inv.DiscountPercent,
			HasDiscount: inv.HasDiscount,
			Available:   inv.Available,
		})
	}

	out := make([]VariantGroup, 0, len(groupOrder))
	for _, attrID := range groupOrder {
		out = append(out, *groups[attrID])
	}
	payload, _ := json.Marshal(variants)
	return out, string(payload)
}

func (p PublicHandler) Shipping(c *gin.Context) {
	customer, ok := helpers.GetAuthUser(c)
	if ok {
		if customer.Cart.CartItem.TotalItemCount <= 0 {
			c.Redirect(http.StatusFound, "/checkout/cart")
			return
		}
	}

	// fee quote from the cart items total (the same rule the order stores)
	quote := p.feeSrv.Quote(c, customer.Cart.CartItem.TotalSalePrice)

	response.CustomerRender(c, http.StatusFound, "shipping",
		gin.H{
			"TITLE":           "اطلاعات ارسال",
			"SHIPPING_FEE":    quote.Shipping,
			"PACKAGING_FEE":   quote.Packaging,
			"SHIPPING_FREE":   quote.Free,
			"FREE_THRESHOLD":  quote.Threshold,
			"GRAND_TOTAL":     quote.Grand,
			"SHIPPING_TITLE":  quote.ShipTitle,
			"PACKAGING_TITLE": quote.PackTitle,
		})
	return
}

func (p PublicHandler) ShowOrderList(c *gin.Context) {

	orderPaginations, err := p.homeSrv.ListOrders(c)
	fmt.Println("----------------------orderPagination:", orderPaginations)
	if err != nil || orderPaginations.Rows == nil {
		//هر خطایی به جز خطای مرتبط با پیدانکردن رکورد اگر وجود داشت اون خطا رو نشون میدیم
		//در غیر این صورت پیغام رکورد یافت نشد به کاربر نشون داده میشه :)
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			response.CustomerRender(c, http.StatusNotFound, "profile_orders", gin.H{
				"TITLE":      "لیست سفارشات",
				"PAGINATION": nil,
				"ACTIVE":     "orders",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"msg": domain_err.SomethingWrongHappened,
			})
			return
		}
	}

	response.CustomerRender(c, http.StatusFound, "profile_orders",
		gin.H{
			"TITLE":          "لیست سفارشات",
			"PAGINATION":     orderPaginations,
			"PrimaryMessage": "لیست سفارشات",
			"ACTIVE":         "orders",
		},
	)
	return
}

func (p PublicHandler) ShowOrderDetails(c *gin.Context) {
	q := c.Param("order_number")
	order, err := p.homeSrv.GetOrderBy(c, q)
	if err != nil || order == nil {
		log.Println("---- [public - handlers]-[ShowOrderDetails]----", err)
		c.JSON(http.StatusOK, gin.H{
			"msg": domain_err.SomethingWrongHappened,
		})
		return
	}
	response.CustomerRender(c, http.StatusFound, "customer_order_details", gin.H{
		"TITLE":  "جزییات سفارش",
		"DATA":   order,
		"ACTIVE": "orders",
	})
	return
}

// SliderCatalog is the «مشاهده همه» page of a product slider: every product
// the slider points at (or its whole category scope) with pagination.
func (p PublicHandler) SliderCatalog(c *gin.Context) {
	slider, err := p.slidersSrv.BySlug(c, c.Param("slug"))
	if err != nil || !slider.IsVisibleNow(time.Now()) {
		response.CustomerRender(c, http.StatusNotFound, "404", gin.H{
			"TITLE": "صفحه یافت نشد",
		})
		return
	}

	page, pageErr := p.slidersSrv.Catalog(c, slider)
	view := responses.ToSlider(slider)

	data := gin.H{
		"TITLE":       view.Title,
		"SLIDER":      view,
		"MEDIA_PATH":  util.GetProductStoragePath(),
		"BANNER_PATH": util.GetBannerStoragePath(),
	}
	if pageErr != nil {
		data["PRIMARY_MESSAGE"] = "محصولی در این بخش یافت نشد"
		data["PAGINATION"] = nil
	} else {
		data["PAGINATION"] = page
	}

	response.CustomerRender(c, http.StatusOK, "catalog", data)
	return
}
