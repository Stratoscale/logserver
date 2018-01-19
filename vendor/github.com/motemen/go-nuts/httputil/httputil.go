package httputil

import (
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
}

func (e *HTTPError) Error() string {
	return e.Status
}

func Successful(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return resp, nil
}

type TransportWrapFunc func(req *http.Request, base http.RoundTripper) (*http.Response, error)

func WrapTransport(base http.RoundTripper, wrapper ...TransportWrapFunc) http.RoundTripper {
	transport := base
	for _, w := range wrapper {
		transport = &wrappedTransport{
			roundTrip: w,
			base:      transport,
		}
	}
	return transport
}

type wrappedTransport struct {
	roundTrip TransportWrapFunc
	base      http.RoundTripper
}

func (t *wrappedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	return t.roundTrip(req, base)
}
