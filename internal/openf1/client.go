package openf1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.openf1.org/v1"

const liveRestrictionMarker = "Live F1 session in progress"

type Session struct {
	SessionKey       int    `json:"session_key"`
	SessionType      string `json:"session_type"`
	SessionName      string `json:"session_name"`
	DateStart        string `json:"date_start"`
	DateEnd          string `json:"date_end"`
	MeetingKey       int    `json:"meeting_key"`
	CircuitShortName string `json:"circuit_short_name"`
	Location         string `json:"location"`
	CountryName      string `json:"country_name"`
	GmtOffset        string `json:"gmt_offset"`
	IsCancelled      bool   `json:"is_cancelled"`
}

type Meeting struct {
	MeetingKey  int    `json:"meeting_key"`
	MeetingName string `json:"meeting_name"`
	Location    string `json:"location"`
	CountryName string `json:"country_name"`
	IsCancelled bool   `json:"is_cancelled"`
}

type Client struct {
	http    *http.Client
	baseURL string
}

type APIError struct {
	URL        string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("fetch %s: status %d: %s", e.URL, e.StatusCode, e.Body)
}

func IsLiveRestriction(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		return false
	}
	return strings.Contains(apiErr.Body, liveRestrictionMarker)
}

func NewClient() *Client {
	return &Client{
		http:    &http.Client{Timeout: 15 * time.Second},
		baseURL: defaultBaseURL,
	}
}

func (c *Client) Sessions(year int) ([]Session, error) {
	var sessions []Session
	if err := c.get(fmt.Sprintf("%s/sessions?year=%d", c.baseURL, year), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (c *Client) Meetings(year int) ([]Meeting, error) {
	var meetings []Meeting
	if err := c.get(fmt.Sprintf("%s/meetings?year=%d", c.baseURL, year), &meetings); err != nil {
		return nil, err
	}
	return meetings, nil
}

func (c *Client) get(url string, dest any) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return &APIError{URL: url, StatusCode: resp.StatusCode, Body: string(body)}
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response from %s: %w", url, err)
	}
	return nil
}
