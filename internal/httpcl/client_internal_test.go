package httpcl

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRetryAfter_Seconds(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": []string{"30"}}}
	if d := parseRetryAfter(resp); d != 30*time.Second {
		t.Fatalf("expected 30s, got %v", d)
	}
}

func TestParseRetryAfter_ZeroSeconds(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": []string{"0"}}}
	if d := parseRetryAfter(resp); d != 0 {
		t.Fatalf("expected 0, got %v", d)
	}
}

func TestParseRetryAfter_HTTPDate_Future(t *testing.T) {
	future := time.Now().Add(5 * time.Second)
	resp := &http.Response{
		Header: http.Header{"Retry-After": []string{future.UTC().Format(http.TimeFormat)}},
	}
	d := parseRetryAfter(resp)
	if d < 3*time.Second || d > 6*time.Second {
		t.Fatalf("expected ~5s from HTTP-date, got %v", d)
	}
}

func TestParseRetryAfter_HTTPDate_Past(t *testing.T) {
	past := time.Now().Add(-5 * time.Second)
	resp := &http.Response{
		Header: http.Header{"Retry-After": []string{past.UTC().Format(http.TimeFormat)}},
	}
	if d := parseRetryAfter(resp); d != 0 {
		t.Fatalf("expected 0 for past date, got %v", d)
	}
}

func TestParseRetryAfter_Missing(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	if d := parseRetryAfter(resp); d != 0 {
		t.Fatalf("expected 0 for missing header, got %v", d)
	}
}

func TestParseRetryAfter_Invalid(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": []string{"garbage"}}}
	if d := parseRetryAfter(resp); d != 0 {
		t.Fatalf("expected 0 for invalid value, got %v", d)
	}
}
