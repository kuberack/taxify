package twilio_client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	twilio "github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"kuberack.com/taxify/internal/models"
)

var twilio_client *TwilioClient

type TwilioClient struct {
	client     *twilio.RestClient
	accountSid string
	serviceId  string
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

func GetTwilioClient() (*TwilioClient, error) {

	// Singleton
	if twilio_client != nil {
		return twilio_client, nil
	}

	// Check if the proxy ip is configured
	// TODO: need to move tclient to the context. Basically each http client for a given
	// external service needs to be available in the context
	var tclient *twilio.RestClient
	var accountSid, authToken, serviceId string
	deploy, exists := os.LookupEnv("TAXIFY_DEPLOY_TYPE")
	if !exists || deploy == "UNIT_TEST" {
		purl, exists := os.LookupEnv("HTTP_PROXY")

		if !exists {
			return &TwilioClient{}, errors.New("proxy not configured for unit test")
		}

		// https://github.com/twilio/twilio-go/blob/main/advanced-examples/custom-http-client.md

		proxyURL, _ := url.Parse(purl)

		// Create your custom Twilio client using the http client and your credentials
		accountSid = "123"
		authToken = "456"
		serviceId = "VAabcd1234abcd1234abcd1234abcd1234" // servicesid must match pattern "^VA[0-9a-fA-F]{32}$"

		twilioHttpClient := &MyClient{
			Client: client.Client{
				Credentials: client.NewCredentials(accountSid, authToken),
			},
			host: proxyURL.Host,
		}
		twilioHttpClient.SetAccountSid(accountSid)
		tclient = twilio.NewRestClientWithParams(twilio.ClientParams{Client: twilioHttpClient})

	} else if deploy == "INTEGRATION_TEST" {
		// Find your Account SID and Auth Token at twilio.com/console
		// and set the environment variables. See http://twil.io/secure
		accountSid, exists = os.LookupEnv("TAXIFY_TWILIO_ACCOUNT_SID")
		if !exists {
			return &TwilioClient{}, errors.New("twilio account sid not present")
		}

		authToken, exists = os.LookupEnv("TAXIFY_TWILIO_AUTH_KEY")
		if !exists {
			return &TwilioClient{}, errors.New("twilio auth key not present")
		}

		// https://console.twilio.com/us1/develop/verify/services
		serviceId, exists = os.LookupEnv("TAXIFY_TWILIO_VERIFY_SERVICE_ID")
		if !exists {
			return &TwilioClient{}, errors.New("twilio verify service id not present")
		}

		tclient = twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSid,
			Password: authToken,
		})
	} else {
		return &TwilioClient{}, errors.New("unknown deployment type")
	}

	twilio_client = &TwilioClient{tclient, accountSid, serviceId}
	return twilio_client, nil
}

func (c *TwilioClient) CreateVerification(phone string) (string, error) {

	// Create a verification
	// https://www.twilio.com/docs/verify/api/verification
	vparams := &verify.CreateVerificationParams{}
	vparams.SetTo(phone)
	vparams.SetChannel("sms")

	resp, err := c.client.VerifyV2.CreateVerification(c.serviceId, vparams)

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return "", errors.New("unable to verify")
	}
	return *resp.Sid, nil
}

func (c *TwilioClient) DoVerificationCheck(record models.User, otp int) error {
	p := &verify.CreateVerificationCheckParams{}
	p.SetTo(record.PhoneNum)
	p.SetCode(strconv.Itoa(otp))

	if resp, err := c.client.VerifyV2.CreateVerificationCheck(record.VerifySid, p); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return errors.New("unable to create check")
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
	return nil
}
