package handlers

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/helpers"
	"shop/internal/pkg/html"
	"shop/internal/pkg/payment/zarinpal"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/sms/kavenegar"
	"strconv"
)

func (p PublicHandler) Payment(c *gin.Context) {
	customer, ok := helpers.GetAuthUser(c)

	//if address fields was empty ,set message
	if !ok || customer.Address.ID == 0 ||
		customer.Address.ReceiverName == "" ||
		customer.Address.ReceiverMobile == "" ||
		customer.Address.ReceiverAddress == "" ||
		customer.Address.ReceiverPostalCode == "" {

		sessions.Set(c, "message", custom_messages.FillAddress)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	zarin, err := zarinpal.NewZarinpal(os.Getenv("ZARINPAL_MERCHANTID"), false)
	if err != nil {
		log.Println("[zarinpal err]:", err)
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, c.Request.Referer())

		return
	}

	_, paymentEntity, inventoryID, paymentError := p.homeSrv.ProcessOrderPayment(c, zarin)

	if paymentError != nil {

		if errors2.Is(paymentError, custom_error.InternalServerErr) {

			sessions.Set(c, "message", custom_error.InternalServerError)
		}
		if errors2.Is(paymentError, custom_error.OutOfStock) {

			errors.Init()
			errors.Add(strconv.Itoa(int(inventoryID)), "موجودی محصول کافی نمی باشد")
			sessions.Set(c, "errors", errors.ToString())
		}

		c.Redirect(http.StatusFound, "/checkout/cart")
		return
	}

	//redirect to bank gateway

	c.Redirect(http.StatusPermanentRedirect, paymentEntity.PaymentURL)
	return
}

func (p PublicHandler) VerifyPayment(c *gin.Context) {

	authority := c.Request.URL.Query().Get("Authority")
	if authority == "" {
		html.CustomerRender(c, http.StatusPermanentRedirect, "404", gin.H{})
		return
	}
	//if status == "NOK" {
	//	//payment was canceled by user of failed
	//}

	//get payment , order(preload OrderItems)
	order, customer, pErr := p.homeSrv.GetPaymentBy(c, authority)

	if pErr != nil || order.Payment.ID <= 0 {
		html.CustomerRender(c, http.StatusPermanentRedirect, "404", gin.H{})
		return
	}

	zarin, err := zarinpal.NewZarinpal(os.Getenv("ZARINPAL_MERCHANTID"), false)
	if err != nil {
		log.Println("[zarinpal err]:", err)
		html.CustomerRender(c, http.StatusInternalServerError, "500", gin.H{})
		return
	}

	// call zarinPal to check payment status
	verified, refID, statusCode, vErr := zarin.PaymentVerification(int(order.Payment.Amount), authority)
	if vErr != nil || !verified || (statusCode != 100 && statusCode != 101) {

	}

	go func() {
		log.Println("-------- call VerifyPayment ----------")
		p.homeSrv.VerifyPayment(c, order, refID, verified)
	}()

	if statusCode == 100 {
		go func() {
			log.Println("send success shopping sms to customer ", customer.Mobile)
			kavenegar.SendSuccShop(customer.Mobile, order.OrderNumber)
		}()
	}

	if verified {
		html.CustomerRender(c, http.StatusOK, "shopping_complete_buy",
			gin.H{
				"ORDER_NUMBER": order.OrderNumber,
			})
		return
	}

	if !verified || vErr != nil {
		//customer canceled payment
		html.CustomerRender(c, http.StatusOK, "shopping_no_complete_buy",
			gin.H{
				"ORDER_NUMBER": order.OrderNumber,
			})
		return
	}

}
