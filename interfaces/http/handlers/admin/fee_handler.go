package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	responses "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/infrastructure/messages"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
)

// feeKind validates the :kind route segment (shipping | packaging).
// An unknown kind falls back to the shipping tariffs.
func feeKind(c *gin.Context) string {
	kind := c.Param("kind")
	if entities.ValidFeeKind(kind) {
		return kind
	}
	return entities.FeeKindShipping
}

// feeID parses the :id route segment (0 = invalid).
func feeID(c *gin.Context) uint {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0
	}
	return uint(id)
}

func (a *AdminHandler) IndexFee(c *gin.Context) {
	kind := feeKind(c)

	rates, err := a.feeSrv.Index(c, kind)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind)
		return
	}

	// highlight the row that applies right now
	var activeID uint
	if active, aErr := a.feeSrv.Active(c, kind); aErr == nil && active != nil {
		activeID = active.ID
	}

	response.Render(c, http.StatusOK, "admin_index_fee", gin.H{
		"TITLE":     "تعرفه‌های " + entities.FeeKindLabel(kind),
		"KIND":      kind,
		"KINDLABEL": entities.FeeKindLabel(kind),
		"FEES":      responses.ToFeeRates(rates, activeID),
	})
	return
}

func (a *AdminHandler) CreateFee(c *gin.Context) {
	kind := feeKind(c)

	response.Render(c, http.StatusOK, "admin_create_fee", gin.H{
		"TITLE":     "تعرفه جدید " + entities.FeeKindLabel(kind),
		"KIND":      kind,
		"KINDLABEL": entities.FeeKindLabel(kind),
		"STATUS":    "1",
	})
	return
}

func (a *AdminHandler) StoreFee(c *gin.Context) {
	kind := feeKind(c)
	_ = c.Request.ParseForm()

	var req requests.CreateFeeRateRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/create")
		return
	}

	if formErrs := requests.ValidateFeeRateForm(c.Request.PostForm, kind); len(formErrs) > 0 {
		flashFormErrors(c, formErrs)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/create")
		return
	}

	if _, err := a.feeSrv.Store(c, req.ToFeeRate(kind)); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/create")
		return
	}

	sessions.Set(c, "message", "تعرفه با موفقیت ایجاد شد")
	c.Redirect(http.StatusFound, "/admins/fees/"+kind)
	return
}

func (a *AdminHandler) EditFee(c *gin.Context) {
	kind := feeKind(c)

	rate, err := a.feeSrv.Show(c, feeID(c))
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind)
		return
	}

	response.Render(c, http.StatusOK, "admin_create_fee", gin.H{
		"TITLE":     "ویرایش تعرفه",
		"KIND":      kind,
		"KINDLABEL": entities.FeeKindLabel(kind),
		"FEE":       responses.ToFeeRate(rate, false),
		"IS_EDIT":   true,
	})
	return
}

func (a *AdminHandler) UpdateFee(c *gin.Context) {
	kind := feeKind(c)
	id := feeID(c)
	if id == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind)
		return
	}
	_ = c.Request.ParseForm()

	var req requests.CreateFeeRateRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/"+c.Param("id")+"/edit")
		return
	}

	if formErrs := requests.ValidateFeeRateForm(c.Request.PostForm, kind); len(formErrs) > 0 {
		flashFormErrors(c, formErrs)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/"+c.Param("id")+"/edit")
		return
	}

	if err := a.feeSrv.Update(c, id, req.ToFeeRate(kind)); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind+"/"+c.Param("id")+"/edit")
		return
	}

	sessions.Set(c, "message", domain_err.SuccessfullyUpdated)
	c.Redirect(http.StatusFound, "/admins/fees/"+kind)
	return
}

func (a *AdminHandler) DeleteFee(c *gin.Context) {
	kind := feeKind(c)
	id := feeID(c)
	if id == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/fees/"+kind)
		return
	}

	if err := a.feeSrv.Delete(c, id); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
	} else {
		sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	}
	c.Redirect(http.StatusFound, "/admins/fees/"+kind)
	return
}
