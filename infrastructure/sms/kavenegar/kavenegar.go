package kavenegar

import (
	"fmt"
	"github.com/kavenegar/kavenegar-go"
	"os"
)

func Send(receptors []string, message string) {

	api := kavenegar.New(os.Getenv("KAVENEGAR_SECRETKEY"))
	sender := ""

	if res, err := api.Message.Send(sender, receptors, message, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Println(err.Error())
		case *kavenegar.HTTPError:
			fmt.Println(err.Error())
		default:
			fmt.Println(err.Error())
		}
	} else {
		for _, r := range res {
			fmt.Println("MessageID 	= ", r.MessageID)
			fmt.Println("Status    	= ", r.Status)
			//...
		}
	}
}

func SendOTP(receptor string, token string) {

	api := kavenegar.New(os.Getenv("KAVENEGAR_SECRETKEY"))
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

type token struct {
}

func SendSuccShop(receptor string, token string) {

	api := kavenegar.New(os.Getenv("KAVENEGAR_SECRETKEY"))
	template := "test-template"

	if res, err := api.Verify.Lookup(receptor, template, token, &kavenegar.VerifyLookupParam{}); err != nil {
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
