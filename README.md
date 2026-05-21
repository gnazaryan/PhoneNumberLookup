# OXIO PhoneNumberLookup
A standalone application exposing a REST/HTTP interface /v1/phone-numbers endpoint that takes request parameters and returns information about a phone number.

## Install & Run

Requires Go 1.22+.

```bash
go mod tidy
go run .
```

Server listens on ':8080'.

Examples:

```bash
curl 'http://localhost:8080/v1/phone-numbers?phoneNumber=%2B12125690123'
# {"phoneNumber":"+12125690123","countryCode":"US","areaCode":"212","localPhoneNumber":"5690123"}

curl 'http://localhost:8080/v1/phone-numbers?phoneNumber=631%20311%208150'
# {"phoneNumber":"6313118150","error":{"countryCode":"required value is missing"}}
```

Run tests:

```bash
go test ./...
```

## Choices

- **Go 1.22 + 'net/http' + 'go-chi/chi/v5'.** chi gives clean routing and a tiny middleware utilities (logger, recoverer, timeout) without bringing in a full framework. Handlers stay plain 'http.HandlerFunc', which keeps 'httptest' easier.

## Deploy to production

The service is stateless and listens on one HTTP port, so deployment is simple.

- Build a static binary with 'CGO_ENABLED=0 go build' and put it in a small distroless Docker image (~12 MB).
- Run it anywhere stateless workloads run — Kubernetes, ECS Fargate, or Cloud Run all work. Two or more replicas behind a load balancer that handles TLS.
- Use '/healthz' for liveness and readiness probes.
- Configure through environment variables (only 'PORT' today).
- Logs go to stdout so the platform's log collector picks them up. Add a Prometheus '/metrics' endpoint before going live.

## Assumptions

- For non space phone numbers I did keep areaCodeLenByDialing, which maps the dialing code to the number of digits for area code.


## Improvements

- Add full list of all countries dialing codes to make the application complete ('dialingToISO' and 'isoToDialing')
- Handle properly the countries which have different number of area code digit lentgths 'areaCodeLenByDialing'