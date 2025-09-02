package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostSignupPhoneIntegration(t *testing.T) {

	url := "/signup/phone"

	// add request body
	number := 9886240527
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
	h := NewServerWithMiddleware()

	// invoke
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
