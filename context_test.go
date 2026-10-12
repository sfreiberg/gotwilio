package gotwilio

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestWithContextCancellation checks that each WithContext method passes its
// context to the HTTP request, so a canceled context stops the request.
func TestWithContextCancellation(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
	}))
	defer srv.Close()

	twilio := NewTwilioClient("AC123", "token")
	twilio.BaseUrl = srv.URL
	twilio.LookupURL = srv.URL

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		call func() error
	}{
		{"SubmitLookupWithContext", func() error {
			_, err := twilio.SubmitLookupWithContext(ctx, LookupReq{PhoneNumber: "+15005550006"})
			return err
		}},
		{"LookupNoCarrierWithContext", func() error {
			_, err := twilio.LookupNoCarrierWithContext(ctx, "+15005550006")
			return err
		}},
		{"GetAvailablePhoneNumbersWithContext", func() error {
			_, _, err := twilio.GetAvailablePhoneNumbersWithContext(ctx, PhoneNumberLocal, "US", AvailablePhoneNumbersOptions{})
			return err
		}},
		{"SendWhatsAppWithContext", func() error {
			_, _, err := twilio.SendWhatsAppWithContext(ctx, "+15005550006", "+15005550006", "hi", "", "")
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); !errors.Is(err, context.Canceled) {
				t.Errorf("err = %v, want context.Canceled", err)
			}
		})
	}

	if n := atomic.LoadInt32(&hits); n != 0 {
		t.Errorf("server received %d requests, want 0", n)
	}
}

// TestWithContextDeadlineDuringBody checks that a deadline hit while reading
// the response body is returned as an error instead of an empty result.
func TestWithContextDeadlineDuringBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send part of the JSON body, then stall until the client gives up.
		w.Write([]byte(`{"available_phone_numbers": [{"phone_number": "+1500`))
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer srv.Close()

	twilio := NewTwilioClient("AC123", "token")
	twilio.BaseUrl = srv.URL
	twilio.LookupURL = srv.URL

	tests := []struct {
		name string
		call func(ctx context.Context) error
	}{
		{"GetAvailablePhoneNumbersWithContext", func(ctx context.Context) error {
			nums, _, err := twilio.GetAvailablePhoneNumbersWithContext(ctx, PhoneNumberLocal, "US", AvailablePhoneNumbersOptions{})
			if err == nil {
				t.Logf("got %d numbers", len(nums))
			}
			return err
		}},
		{"LookupNoCarrierWithContext", func(ctx context.Context) error {
			_, err := twilio.LookupNoCarrierWithContext(ctx, "+15005550006")
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			if err := tt.call(ctx); !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("err = %v, want context.DeadlineExceeded", err)
			}
		})
	}
}
