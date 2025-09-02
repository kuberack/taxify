// This is the old file which uses the non-strict version of the oapi generated code
// The latest file is the server_strict.go

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
	"kuberack.com/taxify/internal/models"
	"kuberack.com/taxify/internal/twilio_client"
)

// optional code omitted

type Server struct {
	tclient *twilio_client.TwilioClient // all clients to appear here
}

func NewServer(twilio_client *twilio_client.TwilioClient) Server {
	return Server{twilio_client}
}

func (Server) GetDriversUserIdVehicles(w http.ResponseWriter, r *http.Request, userId int) {

}

// Signup using OAuth
// (POST /signup/oauth)
func (Server) PostSignupOauth(w http.ResponseWriter, r *http.Request, params PostSignupOauthParams) {

}

// Signup using phone
// (POST /signup/phone)
func (s Server) PostSignupPhone(w http.ResponseWriter, r *http.Request, params PostSignupPhoneParams) {

	// validate input
	switch params.Type {
	case "driver", "rider", "admin":
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user type",
		})
		return
	}

	// get the phone number
	var body PostSignupPhoneJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to decode phone number",
		})
		return
	}

	// Check if phone number is present
	if body.Phone == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "phone number not present",
		})
		return
	}
	// convert the phone number to E.164 format
	var exists bool
	countryCode, exists := os.LookupEnv("TAXIFY_TWILIO_COUNTRY_CODE")
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "twilio country code not present",
		})
		return
	}

	num, err := phonenumbers.Parse(strconv.Itoa(*body.Phone), countryCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error parsing phone number",
		})
		return
	}
	formatted := phonenumbers.Format(num, phonenumbers.E164)

	// Create a verification
	verifySid, err := s.tclient.CreateVerification(formatted)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	// Write an object into the db
	// Write the phone number, verification service id, expiry time, etc. into db
	user := models.User{
		PhoneNum:  formatted,
		VerifySid: verifySid,
	}

	if err := user.Create(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to write to db",
		})
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	// Response
	// TODO: need to validate responses using the openapi
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"userid": user.Id,
	})
}

// Verify using OTP
// (PATCH /signup/phone/{userId}/verify)
func (s Server) PatchSignupPhoneUserIdVerify(w http.ResponseWriter, r *http.Request, userId int) {

	// validate the input userId, lookup the db
	userRecord, err := models.UserByID(userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		fmt.Printf("Bad user id")
		return
	}

	// Get the code from the body
	var body PatchSignupPhoneUserIdVerifyJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to decode OTP",
		})
		return
	}

	// Check if phone number is present
	if body.Otp == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "OTP not present",
		})
		return
	}

	// https://www.twilio.com/docs/verify/api/verification-check
	err = s.tclient.DoVerificationCheck(userRecord, *body.Otp)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to check",
		})
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Printf("phone verification success\n")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": "123",
	})
}
