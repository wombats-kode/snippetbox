package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to out secureheaders
	// middleware, which writes a 200 status code and "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock http handler to our secureHeaders middleware.  Because
	// SecureHeaders *returns* a http.Handler we can call its ServeHTTP() method.
	// passing in the ResponseRecorder and dummy http.request to execute it.
	secureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test
	rs := rr.Result()

	// Check that the middleware has correctly set the Content-Security-Policy
	// header on the response.
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// check that the middleware has correctly set the Referrer-Policy
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check that the middleware has correctly set the X-Content-Type-Options header
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check tjhat the middleware has correctly set the X-Frame-Options header
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check the middleware has correctly set the X-XSS-Protection header
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// Check that the middleware has correclty called the next handler in line
	// and the response status body are as expected
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
