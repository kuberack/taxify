package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"kuberack.com/taxify/internal/models"
)

func TestPostSignupPhoneUnit(t *testing.T) {

	// Get the mock client
	_, mock, err := models.GetDbMockConnection()
	if err != nil {
		t.Errorf("error in dbGet")
	}

	lastInsertId := 101
	sentencePrepare := mock.ExpectPrepare("^insert into users (.+) values (.+)")
	sentencePrepare.ExpectExec().
		WithArgs("+919886240527", "VEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa").
		WillReturnResult(sqlmock.NewResult(int64(lastInsertId), 1))
	columns := []string{"user_id", "phone_number", "verify_sid"}
	rows := sqlmock.NewRows(columns).
		AddRow(lastInsertId, "+919886240527", "VEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	mock.ExpectQuery("^SELECT user_id, phone_number, verify_sid FROM users WHERE user_id = ?").
		WithArgs(lastInsertId).
		WillReturnRows(rows)

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
