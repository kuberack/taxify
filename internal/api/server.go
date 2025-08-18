package api

import (
	"encoding/json"
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

type Server struct{}

func NewServer() Server {
	return Server{}
}

type MyClient struct {
	client.Client
	host string
}

func (c *MyClient) SendRequest(method string, rawURL string, data url.Values, headers map[string]interface{}, body ...byte) (*http.Response, error) {
	// Modify the URL to point to proxy
	if p, err := url.Parse(rawURL); err != nil {
		return nil, err
	} else {
		p.Scheme = "http"
		p.Host = c.host
		rawURL = p.String()
	}

	var resp *http.Response
	var err error
	if resp, err = c.Client.SendRequest(method, rawURL, data, headers, body...); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(resp.StatusCode)
	}
	// Custom code to pre-process response here
	return resp, err
}

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
	user := models.User{
		PhoneNum:  formatted,
		VerifySid: serviceId,
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
func (Server) PatchSignupPhoneUserIdVerify(w http.ResponseWriter, r *http.Request, userId int) {

	// validate the input userId, lookup the db
	userRecord, err := models.UserByID(userId)
	if err != nil {
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
	p.SetTo(userRecord.PhoneNum)
	p.SetCode(strconv.Itoa(*body.Otp))

	if resp, err := tclient.VerifyV2.CreateVerificationCheck(userRecord.VerifySid, p); err != nil {
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
