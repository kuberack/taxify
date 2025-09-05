// This file is currently unused
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
	"kuberack.com/taxify/internal/models"
	"kuberack.com/taxify/internal/twilio_client"
)

// optional code omitted

type TaxifyServer struct {
	tclient *twilio_client.TwilioClient // all clients to appear here
}

var _ StrictServerInterface = (*TaxifyServer)(nil)

func NewServerStrict(twilio_client *twilio_client.TwilioClient) TaxifyServer {
	return TaxifyServer{twilio_client}
}

func (TaxifyServer) GetHealthz(ctx context.Context, request GetHealthzRequestObject) (GetHealthzResponseObject, error) {
	return nil, errors.New("not implemented")
}

func (TaxifyServer) GetDriversUserIdVehicles(ctx context.Context, request GetDriversUserIdVehiclesRequestObject) (GetDriversUserIdVehiclesResponseObject, error) {
	return nil, errors.New("not implemented")

}

// Signup using OAuth
// (POST /signup/oauth)
func (TaxifyServer) PostSignupOauth(ctx context.Context, request PostSignupOauthRequestObject) (PostSignupOauthResponseObject, error) {
	return nil, errors.New("not implemented")
}

// Signup using phone
// (POST /signup/phone)
func (t TaxifyServer) PostSignupPhone(ctx context.Context, request PostSignupPhoneRequestObject) (PostSignupPhoneResponseObject, error) {
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

	// Create a verification
	// https://www.twilio.com/docs/verify/api/verification
	verifySid, err := t.tclient.CreateVerification(formatted)
	if err != nil {
		message := "unable to verify"
		fmt.Printf("error: %s\n", err.Error())
		resp := PostSignupPhone5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	// Write an object into the db
	// Write the phone number, verification service id, expiry time, etc. into db
	user := models.User{
		PhoneNum:  formatted,
		VerifySid: verifySid,
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
func (t TaxifyServer) PatchSignupPhoneUserIdVerify(ctx context.Context, request PatchSignupPhoneUserIdVerifyRequestObject) (PatchSignupPhoneUserIdVerifyResponseObject, error) {

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
	err = t.tclient.DoVerificationCheck(userRecord, *request.Body.Otp)

	if err != nil {
		message := "unable to check"
		fmt.Printf("error: %s\n", err.Error())
		resp := PatchSignupPhoneUserIdVerify5XXJSONResponse{
			Body: struct {
				Message *string "json:\"message,omitempty\""
			}{&message},
			StatusCode: http.StatusInternalServerError}
		return resp, nil
	}

	fmt.Printf("phone verification success\n")
	token := "123"
	return PatchSignupPhoneUserIdVerify200JSONResponse{Token: &token}, nil
}
