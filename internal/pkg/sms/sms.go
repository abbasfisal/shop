package sms

import (
	"fmt"
	"github.com/kavenegar/kavenegar-go"
)

type Service interface {
	Send(receptors []string, message string)
	SendOTP(receptor string, token string)
}

type ServiceManager struct {
	Service Service
}

func (sm *ServiceManager) SetService(service Service) {
	sm.Service = service
}

func (sm *ServiceManager) Send(receptors []string, message string) {
	if sm.Service == nil {
		return
	}
	sm.Service.Send(receptors, message)
}
func (sm *ServiceManager) SendOTP(receptor string, token string) {
	if sm.Service == nil {
		return
	}
	sm.Service.SendOTP(receptor, token)
}

func GetSMSManager() *ServiceManager {
	return &ServiceManager{}
}

func Send(receptors []string, message string) {
	GetSMSManager().Send(receptors, message)
}

func SendOTP(receptor string, token string) {
	GetSMSManager().SendOTP(receptor, token)
}

/**
 *---------------------------
 *			KAVENEGAR
 *---------------------------
 */

type KaveNegar struct {
	ApiKey string
}

func NewKaveNegar(apiKey string) *KaveNegar {
	return &KaveNegar{ApiKey: apiKey}
}

func (k *KaveNegar) Send(receptors []string, message string) {
	api := kavenegar.New(" your apikey ")
	sender := ""

	if res, err := api.Message.Send(sender, receptors, message, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
			break
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
			break
		default:
			fmt.Println(err.Error())
			break
		}
	} else {
		for _, r := range res {
			fmt.Println("MessageID 	= ", r.MessageID)
			fmt.Println("Status    	= ", r.Status)
		}
	}
}

func (k *KaveNegar) SendOTP(receptor string, token string) {

	api := kavenegar.New(k.ApiKey)
	template := "otp-template"
	params := &kavenegar.VerifyLookupParam{}

	if res, err := api.Verify.Lookup(receptor, template, token, params); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
			return
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
			return
		default:
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("MessageID 	= ", res.MessageID)
		fmt.Println("Status    	= ", res.Status)
	}

}
