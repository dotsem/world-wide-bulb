package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type stateResponse struct {
	State      bool  `json:"state"`
	CooldownMs int64 `json:"cooldown_ms"`
	Viewers    int   `json:"viewers"`
}

type historyItem struct {
	ID        int64     `json:"id"`
	State     bool      `json:"state"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type historyResponse struct {
	Items      []historyItem `json:"items"`
	NextBefore int64         `json:"next_before"`
}

type toggleResponse struct {
	ID         string    `json:"id"`
	State      bool      `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	CooldownMs int64     `json:"cooldown_ms"`
}

// region: get_state
func getState(client *http.Client, baseURL string) (*stateResponse, error) {
	resp, err := client.Get(baseURL + "/api/v1/state")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var res stateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	stateStr := "OFF"
	if res.State {
		stateStr = "ON"
	}
	fmt.Printf("State: %s (%d viewers)\n", stateStr, res.Viewers)
	return &res, nil
}

// endregion: get_state

// region: get_history
func getHistory(client *http.Client, baseURL string, limit int) (*historyResponse, error) {
	url := fmt.Sprintf("%s/api/v1/history?limit=%d", baseURL, limit)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var res historyResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	fmt.Printf("History entries: %d\n", len(res.Items))
	return &res, nil
}

// endregion: get_history

// region: post_toggle
func toggleBulb(client *http.Client, baseURL string) (*toggleResponse, error) {
	resp, err := client.Post(baseURL+"/api/v1/toggle", "application/json", strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("toggle failed (%s): %s", resp.Status, string(body))
	}

	var res toggleResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	stateStr := "OFF"
	if res.State {
		stateStr = "ON"
	}
	fmt.Printf("Toggled bulb, ID: %s (State: %s)\n", res.ID, stateStr)
	return &res, nil
}

// endregion: post_toggle

// region: post_reason
func submitReason(client *http.Client, baseURL string, toggleID string, reason string) error {
	payload, err := json.Marshal(map[string]string{
		"id":     toggleID,
		"reason": reason,
	})
	if err != nil {
		return err
	}

	resp, err := client.Post(baseURL+"/api/v1/reason", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("reason submission failed (%s): %s", resp.Status, string(body))
	}

	fmt.Printf("Reason updated for %s\n", toggleID)
	return nil
}

// endregion: post_reason

func main() {
	rawURL := os.Getenv("WWB_URL")
	if rawURL == "" {
		rawURL = "http://localhost:8080"
	}

	baseURL := strings.TrimRight(rawURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1/events")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")

	client := &http.Client{Timeout: 5 * time.Second}
	fmt.Printf("Connecting to REST API at: %s\n", baseURL)

	if _, err := getState(client, baseURL); err != nil {
		fmt.Fprintf(os.Stderr, "error getting state: %v\n", err)
		os.Exit(1)
	}

	if _, err := getHistory(client, baseURL, 10); err != nil {
		fmt.Fprintf(os.Stderr, "error getting history: %v\n", err)
		os.Exit(1)
	}

	toggle, err := toggleBulb(client, baseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error toggling: %v\n", err)
		os.Exit(1)
	}

	if toggle != nil && toggle.ID != "" {
		if err := submitReason(client, baseURL, toggle.ID, "Turning on the bulb for client demo"); err != nil {
			fmt.Fprintf(os.Stderr, "error submitting reason: %v\n", err)
			os.Exit(1)
		}
	}
}
