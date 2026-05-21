package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"PhoneNumberLookup/internal/handler"
)

// A helper function to make an http call used in test cases.
func do(t *testing.T, url string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	handler.PhoneNumbers(w, req)
	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return w.Code, body
}

func TestPhoneNumbers_Success(t *testing.T) {
	code, body := do(t, "/v1/phone-numbers?phoneNumber=%2B12125690123")
	if code != http.StatusOK {
		t.Fatalf("status=%d body=%v", code, body)
	}
	if body["areaCode"] != "212" || body["countryCode"] != "US" {
		t.Fatalf("body=%v", body)
	}
}

func TestPhoneNumbers_MissingCountryCode(t *testing.T) {
	code, body := do(t, "/v1/phone-numbers?phoneNumber=631%20311%208150")
	if code != http.StatusBadRequest {
		t.Fatalf("status=%d", code)
	}
	errObj, _ := body["error"].(map[string]any)
	if errObj["countryCode"] != "required value is missing" {
		t.Fatalf("body=%v", body)
	}
}

func TestPhoneNumbers_InvalidFormat(t *testing.T) {
	code, _ := do(t, "/v1/phone-numbers?phoneNumber=351%2021%20094%202000")
	if code != http.StatusBadRequest {
		t.Fatalf("status=%d", code)
	}
}

func TestPhoneNumbers_MissingPhoneNumber(t *testing.T) {
	code, body := do(t, "/v1/phone-numbers")
	if code != http.StatusBadRequest {
		t.Fatalf("status=%d", code)
	}
	errObj, _ := body["error"].(map[string]any)
	if errObj["phoneNumber"] != "required value is missing" {
		t.Fatalf("body=%v", body)
	}
}
