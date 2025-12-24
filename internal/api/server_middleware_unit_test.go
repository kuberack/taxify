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

	"github.com/DATA-DOG/go-sqlmock"
	"kuberack.com/taxify/internal/models"
)

func TestHealthzUnit(t *testing.T) {

	// Get the mock client
	// This call needs to be done since the health
	// check later will check if the db is created
	_, _, err := models.GetDbMockConnection()
	if err != nil {
		t.Errorf("error in dbGet")
	}

	url := "/healthz"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

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
		t.Errorf("handler returned wrong status code: want %v", http.StatusOK)
	}
}

func TestPostSignupPhoneUnit(t *testing.T) {

	// Get the mock client
	_, mock, err := models.GetDbMockConnection()
	if err != nil {
		t.Errorf("error in dbGet")
	}

	// get the phone number for which we want to test
	phoneNum, ok := os.LookupEnv("TAXIFY_APP_PHONE_NUMBER")
	if !ok {
		t.Errorf("error in phone Number env var")
	}
	countryCode, ok := os.LookupEnv("TAXIFY_APP_COUNTRY_CODE_NUMBER")
	if !ok {
		t.Errorf("error in country code env var")
	}

	lastInsertId := 101
	sentencePrepare := mock.ExpectPrepare("^insert into users (.+) values (.+)")
	sentencePrepare.ExpectExec().
		WithArgs(countryCode+phoneNum, "VEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa").
		WillReturnResult(sqlmock.NewResult(int64(lastInsertId), 1))
	columns := []string{"user_id", "phone_number", "verify_sid"}
	rows := sqlmock.NewRows(columns).
		AddRow(lastInsertId, countryCode+phoneNum, "VEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	mock.ExpectQuery("^SELECT user_id, phone_number, verify_sid FROM users WHERE user_id = ?").
		WithArgs(lastInsertId).
		WillReturnRows(rows)

	url := "/signup/phone"

	// add request body
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
