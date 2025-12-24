package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestPostSignupPhoneIntegration(t *testing.T) {

	url := "/signup/phone"

	// add request body
	// get the phone number for which we want to test
	phoneNum, ok := os.LookupEnv("TAXIFY_APP_PHONE_NUMBER")
	if !ok {
		t.Errorf("error in phone Number env var")
	}
	number, err := strconv.Atoi(phoneNum)
	if err != nil {
		t.Errorf("error in converting str to number")
	}
	data := PostSignupPhoneJSONBody{Phone: &number}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// add query params
	q := req.URL.Query()
	q.Add("type", "driver")
	req.URL.RawQuery = q.Encode() // Encode and assign the query string

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// recorder
	rr := httptest.NewRecorder()

	// create the handler.
	h, err := NewServerWithMiddleware()

	if err != nil {
		log.Fatal(err.Error())
	}

	// invoke
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
