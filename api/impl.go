package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
	twilio "github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// optional code omitted

type Server struct{}

func NewServer() Server {
	return Server{}
}

// in memory db
type inMemoryDBRecord struct {
	phoneNumber string
	verifySid   string
}

var inMemoryDB = make(map[int]*inMemoryDBRecord)
var inMemoryDBRecordId int

var tclient *twilio.RestClient

func (Server) GetDriversUserIdVehicles(w http.ResponseWriter, r *http.Request, userId int) {

}

// Signup using OAuth
// (POST /signup/oauth)
func (Server) PostSignupOauth(w http.ResponseWriter, r *http.Request, params PostSignupOauthParams) {

}

// Signup using phone
// (POST /signup/phone)
func (Server) PostSignupPhone(w http.ResponseWriter, r *http.Request, params PostSignupPhoneParams) {

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

	// Send message to Twilio
	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	accountSid, exists := os.LookupEnv("TAXIFY_TWILIO_ACCOUNT_SID")
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "twilio account sid not present",
		})
		return
	}

	authToken, exists := os.LookupEnv("TAXIFY_TWILIO_AUTH_KEY")
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "twilio auth key not present",
		})
		return
	}

	// https://console.twilio.com/us1/develop/verify/services
	serviceId, exists := os.LookupEnv("TAXIFY_TWILIO_VERIFY_SERVICE_ID")
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "twilio verify service id not present",
		})
		return
	}

	// TODO: need to move it to the context
	tclient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// First, create a verification
	// https://www.twilio.com/docs/verify/api/verification
	vparams := &verify.CreateVerificationParams{}
	vparams.SetTo(formatted)
	vparams.SetChannel("sms")

	if resp, err := tclient.VerifyV2.CreateVerification(serviceId, vparams); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to verify",
		})
		fmt.Printf("error: %s\n", err.Error())
		return
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}

	// Write an object into the db
	// Write the phone number, verification service id, expiry time, etc. into db
	userid := inMemoryDBRecordId
	inMemoryDB[userid] = &inMemoryDBRecord{
		formatted,
		serviceId,
	}
	inMemoryDBRecordId++

	// Response
	// TODO: need to validate responses using the openapi
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"userid": userid,
	})
}

// Verify using OTP
// (PATCH /signup/phone/{userId}/verify)
func (Server) PatchSignupPhoneUserIdVerify(w http.ResponseWriter, r *http.Request, userId int) {

	// validate the input userId, lookup the db
	userRecord, ok := inMemoryDB[userId]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "error",
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
	p := &verify.CreateVerificationCheckParams{}
	p.SetTo(userRecord.phoneNumber)
	p.SetCode(strconv.Itoa(*body.Otp))

	if resp, err := tclient.VerifyV2.CreateVerificationCheck(userRecord.verifySid, p); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to check",
		})
		fmt.Printf("error: %s\n", err.Error())
		return
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
	fmt.Printf("phone verification success\n")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": "123",
	})
}
