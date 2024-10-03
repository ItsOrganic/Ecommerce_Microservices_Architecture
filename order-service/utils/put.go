package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendPutRequest sends a PUT request to the given URL with the provided payload.
// It returns the HTTP response and an error, if any.
func SendPutRequest(url string, payload interface{}) (*http.Response, error) {
	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending PUT request: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the response body for error details
		defer resp.Body.Close()
		return resp, fmt.Errorf("received non-200 response: %s, body: %s", resp.Status, string(body))
	}

	return resp, nil
}
