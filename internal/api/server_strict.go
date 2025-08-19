package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
	twilio "github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"kuberack.com/taxify/internal/models"
)

// optional code omitted

type ServerStrict struct{}

func NewServerStrict() ServerStrict {
	return ServerStrict{}
}

func (ServerStrict) GetDriversUserIdVehicles(ctx context.Context, request GetDriversUserIdVehiclesRequestObject) (GetDriversUserIdVehiclesResponseObject, error) {
	return nil, errors.New("not implemented")

}

// Signup using OAuth
// (POST /signup/oauth)
func (ServerStrict) PostSignupOauth(ctx context.Context, request PostSignupOauthRequestObject) (PostSignupOauthResponseObject, error) {
	return nil, errors.New("not implemented")
}

// Signup using phone
// (POST /signup/phone)
func (ServerStrict) PostSignupPhone(ctx context.Context, request PostSignupPhoneRequestObject) (PostSignupPhoneResponseObject, error) {
	// validate input
	switch request.Params.Type {
	case "driver", "rider", "admin":
	default:
		message := "invalid user type"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusBadRequest}
		return resp, nil
	}

	// Check if phone number is present
	if request.Body.Phone == nil {
		message := "phone number not present"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusBadRequest}
		return resp, nil
	}
	// convert the phone number to E.164 format
	var exists bool
	countryCode, exists := os.LookupEnv("TAXIFY_TWILIO_COUNTRY_CODE")
	if !exists {
		message := "twilio country code not present"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	num, err := phonenumbers.Parse(strconv.Itoa(*request.Body.Phone), countryCode)
	if err != nil {
		message := "Error parsing phone number"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}
	formatted := phonenumbers.Format(num, phonenumbers.E164)

	// Send message to Twilio
	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	accountSid, exists := os.LookupEnv("TAXIFY_TWILIO_ACCOUNT_SID")
	if !exists {
		message := "twilio account sid not present"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	authToken, exists := os.LookupEnv("TAXIFY_TWILIO_AUTH_KEY")
	if !exists {
		message := "twilio auth key not present"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	// https://console.twilio.com/us1/develop/verify/services
	serviceId, exists := os.LookupEnv("TAXIFY_TWILIO_VERIFY_SERVICE_ID")
	if !exists {
		message := "twilio verify service id not present"
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	// Check if the proxy ip is configured
	// TODO: need to move tclient to the context. Basically each http client for a given
	// external service needs to be available in the context
	purl, exists := os.LookupEnv("HTTP_PROXY")
	if !exists {
		tclient = twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSid,
			Password: authToken,
		})
	} else {
		// https://github.com/twilio/twilio-go/blob/main/advanced-examples/custom-http-client.md

		proxyURL, _ := url.Parse(purl)

		// Create your custom Twilio client using the http client and your credentials
		twilioHttpClient := &MyClient{
			Client: client.Client{
				Credentials: client.NewCredentials(accountSid, authToken),
			},
			host: proxyURL.Host,
		}
		twilioHttpClient.SetAccountSid(accountSid)
		tclient = twilio.NewRestClientWithParams(twilio.ClientParams{Client: twilioHttpClient})
	}

	// First, create a verification
	// https://www.twilio.com/docs/verify/api/verification
	vparams := &verify.CreateVerificationParams{}
	vparams.SetTo(formatted)
	vparams.SetChannel("sms")

	if resp, err := tclient.VerifyV2.CreateVerification(serviceId, vparams); err != nil {
		message := "unable to verify"
		fmt.Printf("error: %s\n", err.Error())
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}

	// Write an object into the db
	// Write the phone number, verification service id, expiry time, etc. into db
	user := models.User{
		PhoneNum:  formatted,
		VerifySid: serviceId,
	}
	if err := user.Create(); err != nil {
		message := "unable to write to db"
		fmt.Printf("error: %s\n", err.Error())
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	// Response
	// TODO: need to validate responses using the openapi
	return PostSignupPhone200JSONResponse{Userid: &user.Id}, nil
}

// Verify using OTP
// (PATCH /signup/phone/{userId}/verify)
func (ServerStrict) PatchSignupPhoneUserIdVerify(ctx context.Context, request PatchSignupPhoneUserIdVerifyRequestObject) (PatchSignupPhoneUserIdVerifyResponseObject, error) {

	// validate the input userId, lookup the db
	userRecord, err := models.UserByID(request.UserId)
	if err != nil {
		message := "error"
		fmt.Printf("Bad user id")
		resp := PatchSignupPhoneUserIdVerify5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusBadRequest}
		return resp, nil
	}

	// Check if phone number is present
	if request.Body.Otp == nil {
		message := "OTP not present"
		resp := PatchSignupPhoneUserIdVerify5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusBadRequest}
		return resp, nil
	}

	// https://www.twilio.com/docs/verify/api/verification-check
	p := &verify.CreateVerificationCheckParams{}
	p.SetTo(userRecord.PhoneNum)
	p.SetCode(strconv.Itoa(*request.Body.Otp))

	if resp, err := tclient.VerifyV2.CreateVerificationCheck(userRecord.VerifySid, p); err != nil {
		message := "unable to check"
		fmt.Printf("error: %s\n", err.Error())
		resp := PatchSignupPhoneUserIdVerify5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
	fmt.Printf("phone verification success\n")
	token := "123"
	return PatchSignupPhoneUserIdVerify200JSONResponse{Token: &token}, nil
}
