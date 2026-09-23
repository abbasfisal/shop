package pagination

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	Limit         int
	Page          int
	Sort          string
	TotalRows     int64
	TotalPages    int
	TotalPagesArr []int
	CurrentLink   string
	PrevLink      string
	NextLink      string
	Rows          interface{}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "ID DESC"
	}
	return p.Sort
}

// --------

func Paginate(c *gin.Context, condition string, entity interface{}, pagination *Pagination, db *gorm.DB) (func(db *gorm.DB) *gorm.DB, bool) {

	var totalRows int64
	db.Model(entity).Where(condition).Count(&totalRows)

	//تعداد رکورد هایی که پیدا کرده رو چک میکنیم
	if totalRows <= 0 {
		return nil, false
	}

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	//generate page numbers from 1 to TotalPages
	for i := 1; i <= pagination.TotalPages; i++ {
		pagination.TotalPagesArr = append(pagination.TotalPagesArr, i)
	}

	//CurrentLink
	pagination.CurrentLink = buildLinkCurrent(c)

	// PrevLink (goes to previous page if current page > 1)
	if pagination.Page > 1 {
		pagination.PrevLink = buildLink(c, pagination.Page-1)
	}

	// NextLink (goes to next page if current page < TotalPages)
	if pagination.Page < pagination.TotalPages {
		pagination.NextLink = buildLink(c, pagination.Page+1)
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}, true
}

// Remove 'page' from query and generate PrevLink/NextLink
// و همچنین حفظ کویری های موجود در ادرس برای لینک قبلی ولینک بعدی
func buildLink(c *gin.Context, page int) string {
	// Parse current URL query parameters
	queryParams := c.Request.URL.Query()

	// Remove 'page' if it exists
	queryParams.Del("page")

	// Rebuild the query without 'page' and add the new page number
	queryString := queryParams.Encode()

	// Construct new URL with updated page
	return fmt.Sprintf("?%s&page=%d", queryString, page)
}

// buildLinkCurrent کاربرد این برای حفظ کوپری پارامتر های مرتبط با شماره پیج هست
// یعنی صفحه شماره ۳ باید کوپری پارامترها رو بتونه حفظ کنه و زمانی که روی
// شماره صفحه های دیگه کیک میکنی اون ها رو انتقال بده
func buildLinkCurrent(c *gin.Context) string {
	// Parse current URL query parameters
	queryParams := c.Request.URL.Query()

	// Remove 'page' if it exists
	queryParams.Del("page")

	// Rebuild the query without 'page' and add the new page number
	queryString := queryParams.Encode()

	// Construct new URL with updated page
	return fmt.Sprintf("?%s&page=", queryString)
}
