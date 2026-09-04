package openf1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsLiveRestriction(t *testing.T) {
	live := &APIError{
		URL:        "https://api.openf1.org/v1/meetings?year=2026",
		StatusCode: http.StatusUnauthorized,
		Body:       `{"detail":"Live F1 session in progress. Global API access (including past sessions) is restricted to authenticated users until the session ends."}`,
	}
	if !IsLiveRestriction(live) {
		t.Fatal("expected live restriction")
	}
	if !IsLiveRestriction(fmt.Errorf("wrap: %w", live)) {
		t.Fatal("expected wrapped live restriction")
	}

	other401 := &APIError{StatusCode: http.StatusUnauthorized, Body: `{"detail":"missing api key"}`}
	if IsLiveRestriction(other401) {
		t.Fatal("did not expect a generic 401 to count as live restriction")
	}
	if IsLiveRestriction(fmt.Errorf("fetch failed")) {
		t.Fatal("did not expect a plain error to count as live restriction")
	}
}

func TestClientMeetingsLiveRestriction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/meetings" {
			t.Fatalf("path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"Live F1 session in progress. Global API access is restricted"}`))
	}))
	t.Cleanup(srv.Close)

	client := &Client{http: srv.Client(), baseURL: srv.URL}
	_, err := client.Meetings(2026)
	if !IsLiveRestriction(err) {
		t.Fatalf("got %v", err)
	}
}

func TestClientSessionsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Session{{
			SessionKey:  1,
			SessionName: "Practice 1",
			MeetingKey:  10,
		}})
	}))
	t.Cleanup(srv.Close)

	client := &Client{http: srv.Client(), baseURL: srv.URL}
	sessions, err := client.Sessions(2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].SessionName != "Practice 1" {
		t.Fatalf("got %+v", sessions)
	}
}
