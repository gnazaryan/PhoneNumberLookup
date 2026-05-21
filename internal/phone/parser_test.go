package phone_test

import (
	"errors"
	"testing"

	"PhoneNumberLookup/internal/phone"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		name, in, iso string
		expexcted     phone.ParsedNumber
	}{
		{
			name:      "plus no spaces US/Canada case",
			in:        "+12125690123",
			expexcted: phone.ParsedNumber{PhoneNumber: "+12125690123", CountryCode: "US", AreaCode: "212", LocalPhoneNumber: "5690123"},
		},
		{
			name:      "plus with spaces MX",
			in:        "+52 631 3118150",
			expexcted: phone.ParsedNumber{PhoneNumber: "+526313118150", CountryCode: "MX", AreaCode: "631", LocalPhoneNumber: "3118150"},
		},
		{
			name: "no plus with spaces ES via iso",
			in:   "34 915 872200", iso: "ES",
			expexcted: phone.ParsedNumber{PhoneNumber: "34915872200", CountryCode: "ES", AreaCode: "915", LocalPhoneNumber: "872200"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := phone.Parse(tc.in, tc.iso)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if *actual != tc.expexcted {
				t.Fatalf("actual %+v, want %+v", *actual, tc.expexcted)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []struct {
		name, in, iso string
		expexctedErr  error
	}{
		{"four segments", "351 21 094 2000", "", phone.ErrInvalidFormat},
		{"non digit char", "+1212-5690123", "", phone.ErrInvalidFormat},
		{"empty", "", "", phone.ErrInvalidFormat},
		{"leading space", " +12125690123", "", phone.ErrInvalidFormat},
		{"double space", "+52  631 3118150", "", phone.ErrInvalidFormat},
		{"trailing space", "+12125690123 ", "", phone.ErrInvalidFormat},
		{"missing iso when no plus", "6313118150", "", phone.ErrMissingCountryCode},
		{"unknown iso", "12125690123", "XX", phone.ErrUnknownCountryCode},
		{"iso mismatch with digits", "12125690123", "ES", phone.ErrCountryMismatch},
		{"plus only", "+", "", phone.ErrInvalidFormat},
		{"two segments ambiguous", "+1 2125690123", "", phone.ErrInvalidFormat},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := phone.Parse(tc.in, tc.iso)
			if !errors.Is(err, tc.expexctedErr) {
				t.Fatalf("actual %v, expected %v", err, tc.expexctedErr)
			}
		})
	}
}
