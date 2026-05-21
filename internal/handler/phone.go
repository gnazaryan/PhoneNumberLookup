package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"PhoneNumberLookup/internal/phone"
)

type errorBody struct {
	PhoneNumber string            `json:"phoneNumber,omitempty"`
	Error       map[string]string `json:"error"`
}

// The http handler function for parsing and returning the phon number details.
// PhoneNumbers handles GET /v1/phone-numbers?phoneNumber={phone_number}&countryCode={country_code}
func PhoneNumbers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	raw := q.Get("phoneNumber")
	iso := q.Get("countryCode")

	if raw == "" {
		writeJSON(w, http.StatusBadRequest, errorBody{
			Error: map[string]string{"phoneNumber": "required value is missing"},
		})
		return
	}

	parsed, err := phone.Parse(raw, iso)
	if err != nil {
		writeErr(w, raw, err)
		return
	}
	writeJSON(w, http.StatusOK, parsed)
}

// Construct a user friendly response based on the internal error types.
func writeErr(w http.ResponseWriter, raw string, err error) {
	echo := strings.ReplaceAll(raw, " ", "")
	body := errorBody{PhoneNumber: echo, Error: map[string]string{}}
	switch {
	case errors.Is(err, phone.ErrMissingCountryCode):
		body.Error["countryCode"] = "required value is missing"
	case errors.Is(err, phone.ErrUnknownCountryCode):
		body.Error["countryCode"] = "unknown country code"
	case errors.Is(err, phone.ErrUnknownDialingCode):
		body.Error["phoneNumber"] = "could not identify country from dialing code"
	case errors.Is(err, phone.ErrCountryMismatch):
		body.Error["phoneNumber"] = "does not match provided countryCode"
	default:
		body.Error["phoneNumber"] = "invalid format"
	}
	writeJSON(w, http.StatusBadRequest, body)
}

// Write the josn response output.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
