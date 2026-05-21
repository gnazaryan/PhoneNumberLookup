package phone

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidFormat      = errors.New("invalid phone number format")
	ErrMissingCountryCode = errors.New("country code is required")
	ErrUnknownCountryCode = errors.New("unknown country code")
	ErrUnknownDialingCode = errors.New("could not identify country from dialing code")
	ErrCountryMismatch    = errors.New("phone number does not match provided country code")
)

type ParsedNumber struct {
	PhoneNumber      string `json:"phoneNumber"`
	CountryCode      string `json:"countryCode"`
	AreaCode         string `json:"areaCode"`
	LocalPhoneNumber string `json:"localPhoneNumber"`
}

// Parse the raw inputed number into the parts with the help of the iso country code if provided.
func Parse(raw, isoCountry string) (*ParsedNumber, error) {
	hasPlus, segments, err := tokenize(raw)
	if err != nil {
		return nil, err
	}

	var iso, dialing string
	if hasPlus {
		all := strings.Join(segments, "")
		dialing, iso = lookupByDialing(all)
		if dialing == "" {
			return nil, ErrUnknownDialingCode
		}
	} else {
		if isoCountry == "" {
			return nil, ErrMissingCountryCode
		}
		iso = strings.ToUpper(isoCountry)
		dialing = dialingForISO(iso)
		if dialing == "" {
			return nil, ErrUnknownCountryCode
		}
	}

	digits := strings.Join(segments, "")
	if !strings.HasPrefix(digits, dialing) {
		return nil, ErrCountryMismatch
	}
	rest := digits[len(dialing):]

	area, local, err := splitAreaLocal(segments, dialing, rest, hasPlus)
	if err != nil {
		return nil, err
	}

	display := digits
	if hasPlus {
		display = "+" + digits
	}
	return &ParsedNumber{
		PhoneNumber:      display,
		CountryCode:      iso,
		AreaCode:         area,
		LocalPhoneNumber: local,
	}, nil
}

// Tokenize the raw input and make sure the parts are valid.
func tokenize(raw string) (bool, []string, error) {
	if raw == "" {
		return false, nil, ErrInvalidFormat
	}
	hasPlus := false
	if raw[0] == '+' {
		hasPlus = true
		raw = raw[1:]
	}
	if raw == "" {
		return false, nil, ErrInvalidFormat
	}

	segments := strings.Split(raw, " ")
	if len(segments) > 3 {
		return false, nil, ErrInvalidFormat
	}
	for _, seg := range segments {
		if seg == "" {
			return false, nil, ErrInvalidFormat
		}
		for _, r := range seg {
			if !unicode.IsDigit(r) {
				return false, nil, ErrInvalidFormat
			}
		}
	}
	return hasPlus, segments, nil
}

// Split the area and local parts based on the segments and dial code length.
func splitAreaLocal(segments []string, dialing, rest string, hasPlus bool) (string, string, error) {
	if len(segments) == 3 {
		return segments[1], segments[2], nil
	}
	if len(segments) == 2 {
		return "", "", ErrInvalidFormat
	}

	n, ok := areaCodeLenByDialing[dialing]
	if !ok || n == 0 || len(rest) <= n {
		return "", "", ErrInvalidFormat
	}
	return rest[:n], rest[n:], nil
}
