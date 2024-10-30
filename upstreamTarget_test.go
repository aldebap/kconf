////////////////////////////////////////////////////////////////////////////////
//	upstreamTarget_test.go  -  Oct-29-2024  -  aldebap
//
//	Test cases for Kong Upstream Target Configuration
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test_AddUpstreamTarget unit tests for AddUpstreamTarget() method
func Test_AddUpstreamTarget(t *testing.T) {

	t.Run(">>> AddUpstreamTarget: scenario 1 - error with the request", func(t *testing.T) {

		//	mock for Kong Admin
		var mockKongAdmin *httptest.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer mockKongAdmin.Close()

		//	connect to mock server
		kongServer := NewKongServer(mockKongAdmin.URL, 0)
		if kongServer == nil {
			t.Errorf("fail connectring to mock Kong Admin")
		}

		want := errors.New("fail sending add upstream target command to Kong: 400 Bad Request")
		got := kongServer.AddUpstreamTarget("1234", &KongUpstreamTarget{}, Options{
			verbose:    false,
			jsonOutput: false,
		})

		//	check the invocation result
		if want.Error() != got.Error() {
			t.Errorf("failed checking kong status: error expected: %d result: %d", want, got)
		}
	})

	t.Run(">>> AddUpstreamTarget: scenario 2 - upstreamTarget created successfuly", func(t *testing.T) {

		//	mock for Kong Admin
		var mockKongAdmin *httptest.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"id": "1343894e-404a-4f9e-a982-9e5c0e9d1733",
				"target": "192.168.68.107:8080"
			}`))
		}))
		defer mockKongAdmin.Close()

		//	connect to mock server
		kongServer := NewKongServer(mockKongAdmin.URL, 0)
		if kongServer == nil {
			t.Errorf("fail connectring to mock Kong Admin")
		}

		var want error = nil
		got := kongServer.AddUpstreamTarget("1234", &KongUpstreamTarget{}, Options{
			verbose:    false,
			jsonOutput: false,
		})

		//	check the invocation result
		if want != got {
			t.Errorf("failed checking kong status: success expected: result: %s", got.Error())
		}
	})
}